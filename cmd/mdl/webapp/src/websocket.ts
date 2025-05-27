interface LiveReloadOptions {
	minDelay: number;
	maxDelay: number;
	handshakeTimeout: number;
}

interface LiveReloadMessage {
	command: string;
	protocols?: string[];
	ver?: string;
	path?: string;
}

class Timer {
	private readonly callback: () => void;
	private readonly handler: () => void;
	private running: boolean = false;
	private timeoutId: ReturnType<typeof setTimeout> | null = null;

	constructor(callback: () => void) {
		this.callback = callback;
		this.handler = () => {
			this.running = false;
			this.timeoutId = null;
			this.callback();
		};
	}

	start(timeout: number): void {
		if (this.running) {
			this.stop();
		}
		this.timeoutId = setTimeout(this.handler, timeout);
		this.running = true;
	}

	stop(): void {
		if (this.running && this.timeoutId !== null) {
			clearTimeout(this.timeoutId);
			this.running = false;
			this.timeoutId = null;
		}
	}

	isRunning(): boolean {
		return this.running;
	}
}

/**
 * RefreshConnector implements the livereload protocol to listen for file changes via WebSocket
 */
export class RefreshConnector {
	private static readonly DEFAULT_OPTIONS: LiveReloadOptions = {
		minDelay: 1000,
		maxDelay: 60000,
		handshakeTimeout: 5000
	};

	private static readonly LIVERELOAD_PROTOCOLS = [
		'http://livereload.com/protocols/official-9',
		'http://livereload.com/protocols/2.x-remote-control'
	];

	private readonly uri: string;
	private readonly options: LiveReloadOptions;
	private readonly fileChangeHandler: (path: string) => void;
	
	private socket: WebSocket | null = null;
	private nextDelay: number;
	private connectionDesired: boolean = false;
	private disconnectionReason: string = '';
	
	private readonly handshakeTimeout: Timer;
	private readonly reconnectTimer: Timer;

	constructor(fileChangeHandler: (file: string) => void, options?: Partial<LiveReloadOptions>) {
		this.fileChangeHandler = fileChangeHandler;
		this.options = { ...RefreshConnector.DEFAULT_OPTIONS, ...options };
		this.uri = 'ws://localhost:35729/livereload';
		this.nextDelay = this.options.minDelay;

		this.handshakeTimeout = new Timer(() => this.handleHandshakeTimeout());
		this.reconnectTimer = new Timer(() => this.attemptReconnection());
	}

	connect(): void {
		this.connectionDesired = true;

		if (this.isSocketConnected()) {
			return;
		}

		this.prepareForConnection();
		this.createWebSocket();
	}

	disconnect(): void {
		this.connectionDesired = false;
		this.reconnectTimer.stop();

		if (this.isSocketConnected()) {
			this.disconnectionReason = 'manual';
			this.socket!.close();
		}
	}

	private isSocketConnected(): boolean {
		return this.socket !== null && this.socket.readyState === WebSocket.OPEN;
	}

	private prepareForConnection(): void {
		this.reconnectTimer.stop();
		this.disconnectionReason = 'cannot-connect';
	}

	private createWebSocket(): void {
		this.socket = new WebSocket(this.uri);
		this.socket.onopen = () => this.handleOpen();
		this.socket.onclose = () => this.handleClose();
		this.socket.onmessage = (event) => this.handleMessage(event);
		this.socket.onerror = () => this.handleError();
	}

	private handleOpen(): void {
		this.disconnectionReason = 'handshake-failed';
		this.startHandshake();
	}

	private handleClose(): void {
		console.log(`WebSocket disconnected: ${this.disconnectionReason}. Retry in ${this.nextDelay}ms`);
		this.scheduleReconnection();
	}

	private handleMessage(event: MessageEvent): void {
		try {
			const message: LiveReloadMessage = JSON.parse(event.data);
			this.processMessage(message);
		} catch (error) {
			console.error('Failed to parse WebSocket message:', error);
		}
	}

	private handleError(): void {
		// Error handling is done in onclose
	}

	private processMessage(message: LiveReloadMessage): void {
		switch (message.command) {
			case 'hello':
				this.handleHelloMessage();
				break;
			case 'reload':
				this.handleReloadMessage(message);
				break;
			default:
				console.log('Unknown WebSocket message received:', message);
		}
	}

	private handleHelloMessage(): void {
		this.handshakeTimeout.stop();
		this.nextDelay = this.options.minDelay;
	}

	private handleReloadMessage(message: LiveReloadMessage): void {
		// The livereload server closes connection after sending reload
		// We must reconnect
		this.reconnectTimer.stop();
		this.connect();

		if (message.path) {
			this.fileChangeHandler(message.path);
		}
	}

	private startHandshake(): void {
		const helloMessage: LiveReloadMessage = {
			command: 'hello',
			protocols: RefreshConnector.LIVERELOAD_PROTOCOLS,
			ver: '3.3.1'
		};
		
		this.sendCommand(helloMessage);
		this.handshakeTimeout.start(this.options.handshakeTimeout);
	}

	private handleHandshakeTimeout(): void {
		if (this.isSocketConnected()) {
			this.disconnectionReason = 'handshake-timeout';
			this.socket!.close();
		}
	}

	private attemptReconnection(): void {
		if (this.connectionDesired) {
			this.connect();
		}
	}

	private scheduleReconnection(): void {
		if (!this.connectionDesired) {
			return; // Don't reconnect after manual disconnection
		}

		if (!this.reconnectTimer.isRunning()) {
			this.reconnectTimer.start(this.nextDelay);
			this.nextDelay = Math.min(this.options.maxDelay, this.nextDelay * 2);
		}
	}

	private sendCommand(command: LiveReloadMessage): void {
		if (this.isSocketConnected()) {
			this.socket!.send(JSON.stringify(command));
		}
	}
}