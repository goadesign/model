/**
 * Undo functionality
 * at every change in the document, Undo can save a new version
 * so the user can "undo" and "redo" changes by reverting to an
 * older version of the document
 */

export class Undo<Doc> {
	private readonly versions: Doc[] = [];
	private pos: number = 0;
	private lastSavedPos: number = 0;
	private readonly exportDoc: () => Doc;
	private readonly importDoc: (d: Doc) => void;
	change: () => void;
	private tmpPreviousState: Doc | null = null;

	constructor(id: string, exportDoc: () => Doc, importDoc: (d: Doc) => void) {
		this.exportDoc = exportDoc;
		this.importDoc = importDoc;
		this.change = debounce(() => this.saveNow(), 300);
	}

	// Store the state previous to the changes collected in the debounce period
	beforeChange() {
		if (!this.tmpPreviousState) {
			this.tmpPreviousState = this.deepClone(this.exportDoc());
		}
	}

	length() {
		return this.versions.length;
	}

	currentState() {
		return this.deepClone(this.versions[this.pos - 1]);
	}

	private saveNow() {
		if (!this.tmpPreviousState) {
			throw Error("undo.change() was called without previously calling undo.beforeChange()!");
		}
		
		this.versions[this.pos] = this.deepClone(this.exportDoc());
		this.versions[this.pos - 1] = this.tmpPreviousState;
		this.tmpPreviousState = null;
		this.pos += 1;
		
		// Remove anything that might be on top of this version
		this.versions.splice(this.pos);
	}

	private deepClone(doc: Doc): Doc {
		// Use modern structuredClone if available, fallback to JSON
		if (typeof structuredClone !== 'undefined') {
			return structuredClone(doc);
		}
		return JSON.parse(JSON.stringify(doc));
	}

	undo() {
		if (this.pos < 2) return;
		this.pos -= 1;
		const doc = this.versions[this.pos - 1];
		this.importDoc(this.deepClone(doc));
	}

	redo() {
		if (this.pos > this.versions.length - 1) return;
		const doc = this.versions[this.pos];
		this.importDoc(this.deepClone(doc));
		this.pos += 1;
	}

	changed() {
		return this.pos !== this.lastSavedPos;
	}

	setSaved() {
		this.lastSavedPos = this.pos;
	}
}

function debounce(func: () => void, wait: number) {
	let timeout: ReturnType<typeof setTimeout>;
	return function () {
		const context = this;
		const later = function () {
			timeout = null;
			func.apply(context);
		};
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
	};
}