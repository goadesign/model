/**
 * Undo functionality
 * at every change in the document, Undo can save a new version
 * so the user can "undo" and "redo" changes by reverting to an
 * older version of the document
 */
export class Undo<Doc> {

	private readonly versions: Doc[]
	private pos: number
	private readonly exportDoc: () => Doc
	private readonly importDoc: (d: Doc) => void
	change: () => void

	constructor(exportDoc: () => Doc, importDoc: (d: Doc) => void) {
		this.versions = []
		this.exportDoc = exportDoc
		this.importDoc = importDoc
		this.pos = 0
		this.change = debounce(this.saveNow, 500)
	}

	private saveNow() {
		this.versions[this.pos] = this.deepClone(this.exportDoc())
		this.pos += 1
		// remove anything that might be on top of this version
		this.versions.splice(this.pos, this.versions.length - this.pos)
	}

	private deepClone(doc: Doc): Doc {
		return JSON.parse(JSON.stringify(doc))
	}

	undo() {
		if (this.pos < 2) return
		this.pos -= 1
		const doc = this.versions[this.pos - 1]
		this.importDoc(this.deepClone(doc))
	}

	redo() {
		if (this.pos > this.versions.length - 1) return
		const doc = this.versions[this.pos]
		this.importDoc(this.deepClone(doc))
		this.pos += 1
	}
}

function debounce(func: () => void, wait: number) {
	let timeout: number;
	return function () {
		let context = this;
		let later = function () {
			timeout = null;
			func.apply(context);
		};
		clearTimeout(timeout);
		timeout = setTimeout(later, wait);
	}
}