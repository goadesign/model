import React, {FC} from "react";

interface Combination {
	ctrl?: boolean;
	shift?: boolean;
	alt?: boolean;
	wheel?: boolean
	key?: string;
	click?: boolean;
}

interface Shortcut {
	id: string,
	help: string,
	combinations: Combination[]
}

export const UNDO = 'undo'
export const REDO = 'redo'
export const ADD_VERTEX = 'add-vertex'
export const ADD_LABEL_VERTEX = 'add-label-vertex'
export const DEL_VERTEX = 'del-vertex'

export const ZOOM_IN = 'zoom-in'
export const ZOOM_OUT = 'zoom-out'
export const ZOOM_FIT = 'zoom-fit'
export const ZOOM_100 = 'zoom-100'

export const SELECT_ALL = 'select-all'
export const DESELECT = 'deselect'

export const MOVE_LEFT = 'move-left'
export const MOVE_RIGHT = 'move-right'
export const MOVE_UP = 'move-up'
export const MOVE_DOWN = 'move-down'
export const MOVE_LEFT_FAST = 'move-left-fast'
export const MOVE_RIGHT_FAST = 'move-right-fast'
export const MOVE_UP_FAST = 'move-up-fast'
export const MOVE_DOWN_FAST = 'move-down-fast'

const shortcuts: { name: string; list: Shortcut[] }[] = [
	{
		name: 'History',
		list: [
			{
				id: UNDO,
				help: 'Undo',
				combinations: [
					{ctrl: true, key: 'z'},
				]
			},
			{
				id: REDO,
				help: 'Redo',
				combinations: [
					{ctrl: true, shift: true, key: 'z'},
					{ctrl: true, key: 'y'},
				]
			}

		],
	},
	{
		name: 'Relationship editing',
		list: [
			{
				id: ADD_VERTEX,
				help: 'Add relationship vertex',
				combinations: [
					{alt: true, click: true},
				]
			},
			{
				id: ADD_LABEL_VERTEX,
				help: 'Add label anchor relationship vertex',
				combinations: [
					{alt: true, shift: true, click: true},
				]
			},
			{
				id: DEL_VERTEX,
				help: 'Remove relationship vertex',
				combinations: [
					{key: 'DELETE'},
					{key: 'BACKSPACE'}
				]
			},
		]
	},
	{
		name: 'Zoom',
		list: [
			{
				id: ZOOM_IN,
				help: 'Zoom in',
				combinations: [
					{ctrl: true, key: '='},
					{ctrl: true, wheel: true}
				]
			},
			{
				id: ZOOM_OUT,
				help: 'Zoom out',
				combinations: [
					{ctrl: true, key: '-'},
					{ctrl: true, wheel: true}
				]
			},
			{
				id: ZOOM_FIT,
				help: 'Zoom - fit',
				combinations: [{ctrl: true, key: '9'}]
			},
			{
				id: ZOOM_100,
				help: 'Zoom 100%',
				combinations: [{ctrl: true, key: '0'}]
			}
		]
	},
	{
		name: 'Select',
		list: [
			{
				id: SELECT_ALL,
				help: 'Select All',
				combinations: [{ctrl: true, key: 'a'}]
			},
			{
				id: DESELECT,
				help: 'Deselect',
				combinations: [{key: 'ESC'}]
			}
		]
	},
	{
		name: 'Move',
		list: [
			{
				id: MOVE_UP,
				help: 'Move up',
				combinations: [{key: 'UP'}]
			},
			{
				id: MOVE_UP_FAST,
				help: 'Move up fast',
				combinations: [{key: 'UP', shift: true}]
			},
			{
				id: MOVE_RIGHT,
				help: 'Move right',
				combinations: [{key: 'RIGHT'}]
			},
			{
				id: MOVE_RIGHT_FAST,
				help: 'Move right fast',
				combinations: [{key: 'RIGHT', shift: true}]
			},
			{
				id: MOVE_DOWN,
				help: 'Move down',
				combinations: [{key: 'DOWN'}]
			},
			{
				id: MOVE_DOWN_FAST,
				help: 'Move down fast',
				combinations: [{key: 'DOWN', shift: true}]
			},
			{
				id: MOVE_LEFT,
				help: 'Move left',
				combinations: [{key: 'LEFT'}]
			},
			{
				id: MOVE_LEFT_FAST,
				help: 'Move left fast',
				combinations: [{key: 'LEFT', shift: true}]
			},
		]
	}
]

const shortcutMap = shortcuts
	.reduce((lst, s) => lst.concat(s.list), [] as Shortcut[])
	.reduce<{ [k: string]: Shortcut }>((map, s) => {
		map[s.id] = s;
		return map
	}, {})

const checkKey = (e: KeyboardEvent | MouseEvent, shortcut: Shortcut, click: boolean, wheel: boolean) => {
	return shortcut.combinations.some(c => {
		if (Boolean(c.shift) != e.shiftKey) return false
		if (Boolean(c.ctrl) != e.ctrlKey) return false
		if (Boolean(c.alt) != e.altKey) return false
		if (click) return c.click
		if (wheel) return c.wheel
		if (c.key) {
			const ke = e as KeyboardEvent
			if (c.key == 'DELETE') return ke.key == 'Delete'
			if (c.key == 'BACKSPACE') return ke.key == 'Backspace'
			if (c.key == 'ESC') return ke.key == 'Escape'
			if (c.key == 'UP') return ke.key == 'ArrowUp'
			if (c.key == 'DOWN') return ke.key == 'ArrowDown'
			if (c.key == 'LEFT') return ke.key == 'ArrowLeft'
			if (c.key == 'RIGHT') return ke.key == 'ArrowRight'
			return c.key.toLowerCase() == ke.key.toLowerCase()
		}
		return false
	})
}

export const findShortcut = (e: KeyboardEvent | MouseEvent, click = false, wheel = false) => {
	return Object.keys(shortcutMap).find(k => checkKey(e, shortcutMap[k], click, wheel) && k)
}

const comboText = (c: Combination) => {
	return [
		c.ctrl && 'CTRL',
		c.shift && 'SHIFT',
		c.alt && 'ALT',
		c.key && (c.key.length > 1 ? c.key : `"${c.key.toUpperCase()}"`),
		c.click && 'CLICK',
		c.wheel && 'WHEEL'
	].filter(Boolean).join(' + ')
}

export const Help: FC = () => {
	return <div className="popover">
		<h1>Shortcuts</h1>
		<table>
			<tbody>
			{
				shortcuts.map(section => <>
					<tr>
						<th colSpan={2}>{section.name}</th>
					</tr>
					{section.list.map(item => <tr>
						<td>{item.combinations.map(comboText).join(', ')}</td>
						<td>{item.help}</td>
					</tr>)}
				</>)
			}
			</tbody>
		</table>
	</div>
}
