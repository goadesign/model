class Timer {
	private readonly func: () => void
	private readonly _handler: () => void
	private running: boolean
	private id: number

	constructor(func: () => void) {
		this.func = func;
		this.running = false;
		this.id = null;

		this._handler = () => {
			this.running = false;
			this.id = null;
			this.func();
		};
	}

	start(timeout: number) {
		if (this.running) {
			clearTimeout(this.id);
		}
		this.id = setTimeout(this._handler, timeout);
		this.running = true;
	}

	stop() {
		if (this.running) {
			clearTimeout(this.id);
			this.running = false;
			this.id = null;
		}
	}

	isRunning() {
		return this.running
	}
}

const options = {
	minDelay: 1000,
	maxDelay: 60000,
	handshakeTimeout: 5000
};

/**
 * RefreshConnector implements livereload protocol to listen for changed files via websocket.
 * 
 */
export class RefreshConnector {
	private readonly _uri: string
	private socket: WebSocket;
	private _nextDelay: number;
	private _connectionDesired: boolean;
	private _handshakeTimeout: Timer;
	private _disconnectionReason: string;
	private _reconnectTimer: Timer;
	private handler: (path: string) => void;

	constructor(handler: (file: string) => void) {
		this.handler = handler

		this._uri = `ws://localhost:35729/livereload`;

		this._nextDelay = options.minDelay;
		this._connectionDesired = false;


		this._handshakeTimeout = new Timer(() => {
			if (!this._isSocketConnected()) {
				return;
			}
			this._disconnectionReason = 'handshake-timeout';
			return this.socket.close();
		});

		this._reconnectTimer = new Timer(() => {
			if (!this._connectionDesired) {
				// shouldn't hit this, but just in case
				return;
			}
			return this.connect();
		});

		this.connect();
	}

	_isSocketConnected() {
		return this.socket && (this.socket.readyState === WebSocket.OPEN);
	}

	connect() {
		this._connectionDesired = true;

		if (this._isSocketConnected()) {
			return;
		}

		// prepare for a new connection
		this._reconnectTimer.stop();
		this._disconnectionReason = 'cannot-connect';

		this.socket = new WebSocket(this._uri);
		this.socket.onopen = e => this._onopen();
		this.socket.onclose = () => this._onclose();
		this.socket.onmessage = e => this._onmessage(e);
		this.socket.onerror = e => null;
	}

	disconnect() {
		this._connectionDesired = false;
		this._reconnectTimer.stop(); // in case it was running

		if (!this._isSocketConnected()) {
			return;
		}
		this._disconnectionReason = 'manual';
		return this.socket.close();
	}

	_scheduleReconnection() {
		if (!this._connectionDesired) {
			// don't reconnect after manual disconnection
			return;
		}
		if (!this._reconnectTimer.isRunning()) {
			this._reconnectTimer.start(this._nextDelay);
			this._nextDelay = Math.min(options.maxDelay, this._nextDelay * 2);
		}
	}

	_sendCommand(command: any) {
		if (this._isSocketConnected())
			this.socket.send(JSON.stringify(command));
	}

	_closeOnError() {
		this._handshakeTimeout.stop();
		this._disconnectionReason = 'error';
		return this.socket.close();
	}

	_onopen() {
		// console.log('WS connected')
		this._disconnectionReason = 'handshake-failed';

		// start handshake
		const protocols = [
			'http://livereload.com/protocols/official-9',
			'http://livereload.com/protocols/2.x-remote-control'
		];
		const hello = {command: 'hello', protocols: protocols, ver: '3.3.1'};
		this._sendCommand(hello);
		return this._handshakeTimeout.start(options.handshakeTimeout);
	}

	_onclose() {
		console.log(`WS disconnected: ${this._disconnectionReason}. Retry in ${this._nextDelay}`);
		return this._scheduleReconnection();
	}

	_onmessage(e: WebSocketMessageEvent) {
		const msg = JSON.parse(e.data)
		if (msg && msg.command == 'hello') {
			this._handshakeTimeout.stop();
			this._nextDelay = options.minDelay;
		} else if (msg.command == 'reload') {
			// the livereload server stops connection after sending that. we must connect back
			this._reconnectTimer.stop()
			this.connect()

			this.handler.apply(null, [msg.path])
		} else {
			console.log('WS unknown message received:', msg);
		}
	}
}