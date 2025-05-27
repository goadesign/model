import React, {FC} from "react";
import { isMac, getModifierKeyName, getModifierKeyProperty } from './utils/platform';

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

export const SAVE = 'save'

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
export const MOVE_LEFT_FINE = 'move-left-fine'
export const MOVE_RIGHT_FINE = 'move-right-fine'
export const MOVE_UP_FINE = 'move-up-fine'
export const MOVE_DOWN_FINE = 'move-down-fine'

export const PAN_VIEW = 'pan-view'
export const SELECT_ELEMENT = 'select-element'
export const MULTI_SELECT = 'multi-select'
export const BOX_SELECT = 'box-select'
export const MOVE_ELEMENTS = 'move-elements'

export const HELP = 'help'

// New shortcuts for toolbar buttons
export const TOGGLE_DRAG_MODE = 'toggle_drag_mode'
export const ALIGN_HORIZONTAL = 'align_horizontal'
export const ALIGN_VERTICAL = 'align_vertical'
export const DISTRIBUTE_HORIZONTAL = 'distribute_horizontal'
export const DISTRIBUTE_VERTICAL = 'distribute_vertical'
export const AUTO_LAYOUT = 'auto_layout'
export const RESET_POSITION = 'reset_position'
export const TOGGLE_GRID = 'toggle_grid'
export const TOGGLE_SNAP_TO_GRID = 'toggle_snap_to_grid'
export const SNAP_ALL_TO_GRID = 'snap_all_to_grid'

const shortcuts: { name: string; list: Shortcut[] }[] = [
	{
		name: 'Help',
		list: [
			{
				id: HELP,
				help: 'Show/hide this help',
				combinations: [
					{key: '?', shift: true},
					{key: 'F1', shift: true}
				]
			}
		]
	},
	{
		name: 'File',
		list: [
			{
				id: SAVE,
				help: 'Save',
				combinations: [{key: 's', ctrl: true}]
			}
		]
	},
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
					{ctrl: true, key: '='}
				]
			},
			{
				id: ZOOM_OUT,
				help: 'Zoom out',
				combinations: [
					{ctrl: true, key: '-'}
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
			},
			{
				id: 'wheel_zoom',
				help: 'Zoom in/out with mouse wheel',
				combinations: [{wheel: true}]
			}
		]
	},
	{
		name: 'Mouse Interactions',
		list: [
			{
				id: PAN_VIEW,
				help: 'Pan view (drag empty space)',
				combinations: [{click: true}]
			},
			{
				id: SELECT_ELEMENT,
				help: 'Select element',
				combinations: [{click: true}]
			},
			{
				id: MULTI_SELECT,
				help: 'Add/remove from selection',
				combinations: [{shift: true, click: true}]
			},
			{
				id: BOX_SELECT,
				help: 'Box selection (drag empty space)',
				combinations: [{shift: true, click: true}]
			},
			{
				id: MOVE_ELEMENTS,
				help: 'Move selected elements',
				combinations: [{click: true}]
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
				help: 'Move up (grid increment)',
				combinations: [{key: 'UP'}]
			},
			{
				id: MOVE_UP_FINE,
				help: 'Move up (1 pixel)',
				combinations: [{key: 'UP', shift: true}]
			},
			{
				id: MOVE_RIGHT,
				help: 'Move right (grid increment)',
				combinations: [{key: 'RIGHT'}]
			},
			{
				id: MOVE_RIGHT_FINE,
				help: 'Move right (1 pixel)',
				combinations: [{key: 'RIGHT', shift: true}]
			},
			{
				id: MOVE_DOWN,
				help: 'Move down (grid increment)',
				combinations: [{key: 'DOWN'}]
			},
			{
				id: MOVE_DOWN_FINE,
				help: 'Move down (1 pixel)',
				combinations: [{key: 'DOWN', shift: true}]
			},
			{
				id: MOVE_LEFT,
				help: 'Move left (grid increment)',
				combinations: [{key: 'LEFT'}]
			},
			{
				id: MOVE_LEFT_FINE,
				help: 'Move left (1 pixel)',
				combinations: [{key: 'LEFT', shift: true}]
			},
		]
	},
			{
			name: 'View',
			list: [
				{
					id: TOGGLE_DRAG_MODE,
					help: 'Toggle between pan and select mode',
					combinations: [{key: 't'}]
				},
				{
					id: RESET_POSITION,
					help: 'Reset position and view',
					combinations: [{key: 'Home', ctrl: true}]
				}
			]
		},
	{
		name: 'Alignment',
		list: [
			{
				id: ALIGN_HORIZONTAL,
				help: 'Align selected elements horizontally',
				combinations: [{key: 'h', ctrl: true, shift: true}]
			},
			{
				id: ALIGN_VERTICAL,
				help: 'Align selected elements vertically',
				combinations: [{key: 'a', ctrl: true, shift: true}]
			},
			{
				id: DISTRIBUTE_HORIZONTAL,
				help: 'Distribute selected elements horizontally',
				combinations: [{key: 'h', ctrl: true, alt: true}]
			},
			{
				id: DISTRIBUTE_VERTICAL,
				help: 'Distribute selected elements vertically',
				combinations: [{key: 'v', ctrl: true, alt: true}]
			}
		]
	},
	{
		name: 'Layout',
		list: [
			{
				id: AUTO_LAYOUT,
				help: 'Auto layout all elements',
				combinations: [{key: 'l', ctrl: true}]
			}
		]
	},
	{
		name: 'Grid',
		list: [
			{
				id: TOGGLE_GRID,
				help: 'Toggle grid visibility',
				combinations: [{key: 'g', ctrl: true}]
			},
			{
				id: TOGGLE_SNAP_TO_GRID,
				help: 'Toggle snap to grid',
				combinations: [{key: 'g', ctrl: true, shift: true}]
			},
			{
				id: SNAP_ALL_TO_GRID,
				help: 'Snap all elements to grid',
				combinations: [{key: 'g', ctrl: true, alt: true}]
			}
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
		// Use platform-appropriate modifier key
		if (c.ctrl) {
			const modifierKey = getModifierKeyProperty(e as KeyboardEvent);
			if (!modifierKey) return false;
		}
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
			return c.key && ke.key && c.key.toLowerCase() == ke.key.toLowerCase()
		}
		return false
	})
}

export const findShortcut = (e: KeyboardEvent | MouseEvent, click = false, wheel = false) => {
	// Find all matching shortcuts
	const matches = Object.keys(shortcutMap).filter(k => checkKey(e, shortcutMap[k], click, wheel))
	
	if (matches.length === 0) return undefined
	if (matches.length === 1) return matches[0]
	
	// If multiple matches, prefer the one with more specific modifiers
	// Sort by number of modifiers (shift, ctrl, alt) in descending order
	const sortedMatches = matches.sort((a, b) => {
		const aShortcut = shortcutMap[a]
		const bShortcut = shortcutMap[b]
		
		const aModifiers = aShortcut.combinations[0]
		const bModifiers = bShortcut.combinations[0]
		
		const aCount = (aModifiers.shift ? 1 : 0) + (aModifiers.ctrl ? 1 : 0) + (aModifiers.alt ? 1 : 0)
		const bCount = (bModifiers.shift ? 1 : 0) + (bModifiers.ctrl ? 1 : 0) + (bModifiers.alt ? 1 : 0)
		
		return bCount - aCount // Descending order (more modifiers first)
	})
	
	return sortedMatches[0]
}

const comboText = (c: Combination) => {
	return [
		c.ctrl && getModifierKeyName().toUpperCase(),
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
