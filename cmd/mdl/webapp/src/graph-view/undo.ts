/**
 * Undo functionality
 * at every change in the document, Undo can save a new version
 * so the user can "undo" and "redo" changes by reverting to an
 * older version of the document
 */

// local memory cache
const cache = new Map<string, { versions: any; pos: number; lastSavedPos: number }>();


export class Undo<Doc> {

	private readonly id: string
	private readonly versions: Doc[]
	private pos: number
	private lastSavedPos: number
	private readonly exportDoc: () => Doc
	private readonly importDoc: (d: Doc) => void
	change: () => void
	private tmpPreviousState: Doc

	constructor(id: string, exportDoc: () => Doc, importDoc: (d: Doc) => void) {
		this.exportDoc = exportDoc
		this.importDoc = importDoc
		this.change = debounce(this.saveNow, 300)

		this.id = id

		if (cache.has(this.id)) {
			const c = cache.get(this.id)
			this.versions = c.versions
			this.pos = c.pos
			this.lastSavedPos = c.lastSavedPos
		} else {
			this.pos = 1
			this.lastSavedPos = 1
			this.versions = []
		}
	}

	// the state previous to the changes collected in the debounce period is stored here
	// we use this as the state to revert to on "Undo"
	// this state might differ from the previously saved one because we allow model changes
	// via websocket
	beforeChange() {
		this.tmpPreviousState || (this.tmpPreviousState = this.deepClone(this.exportDoc()))
	}

	length() {
		return this.versions.length
	}

	currentState() {
		return this.deepClone(this.versions[this.pos - 1])
	}

	private saveNow() {
		if (!this.tmpPreviousState) throw Error("undo.change() was called without previously calling undo.beforeChange()!")
		this.versions[this.pos] = this.deepClone(this.exportDoc())
		this.versions[this.pos - 1] = this.tmpPreviousState
		this.tmpPreviousState = null
		this.pos += 1
		// remove anything that might be on top of this version
		this.versions.splice(this.pos, this.versions.length - this.pos)
		this.saveCache()
	}

	private saveCache() {
		cache.set(this.id, {
			versions: this.versions,
			pos: this.pos,
			lastSavedPos: this.lastSavedPos
		})
	}

	private deepClone(doc: Doc): Doc {
		return JSON.parse(JSON.stringify(doc))
	}

	undo() {
		if (this.pos < 2) return
		this.pos -= 1
		const doc = this.versions[this.pos - 1]
		this.importDoc(this.deepClone(doc))
		this.saveCache()
	}

	redo() {
		if (this.pos > this.versions.length - 1) return
		const doc = this.versions[this.pos]
		this.importDoc(this.deepClone(doc))
		this.pos += 1
		this.saveCache()
	}

	changed() {
		return this.pos != this.lastSavedPos
	}

	setSaved() {
		this.lastSavedPos = this.pos
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