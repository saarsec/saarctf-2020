import { VisitedField } from './Cave';

export const CAVE_ROCK = 0;
export const CAVE_FREE = 1;
export const CAVE_START = 2;

const DRAWSIZE = 3;


function readInt(buffer, pos) {
	return buffer[pos] | buffer[pos + 1] << 8 | buffer[pos + 2] << 16 | buffer[pos + 3] << 24;
}

export default class ParsedCave {
	public width: number;
	public height: number;
	public data: number[];
	public startpoint: {row: number, col: number, index?: number};

	constructor(description) {
		if (typeof description === 'object' && description.width > 0) {
			this.width = description.width;
			this.height = description.height;
			this.data = new Array(this.width * this.height);
			for (let i = 0; i < this.width * this.height; i++) {
				switch (description.data.charAt(i)) {
					case ' ':
						this.data[i] = CAVE_FREE;
						break;
					case 'S':
						this.data[i] = CAVE_START;
						this.startpoint = {index: i, row: Math.floor(i / this.width), col: i % this.width};
						break;
					default:
						this.data[i] = CAVE_ROCK;
				}
			}
		} else {
			if (description instanceof ArrayBuffer) {
				description = new Uint8Array(description);
			}
			console.log(description);
			this.width = readInt(description, 0);
			this.height = readInt(description, 4);
			this.startpoint = {row: readInt(description, 8), col: readInt(description, 12)};
			console.log(this.width, this.height, this.startpoint);
			this.data = new Array(this.width * this.height);
			for (let i = 0; i < this.width * this.height; i++) {
				let bit = (description[16 + Math.floor(i / 8)] >> (i % 8)) & 1;
				if (bit) {
					this.data[i] = CAVE_FREE;
				} else {
					this.data[i] = CAVE_ROCK;
				}
			}
		}
	}


	drawToCanvas(canvas, path: VisitedField[]) {
		canvas.width = DRAWSIZE * (this.width + 2);
		canvas.height = DRAWSIZE * (this.height + 2);
		let ctx = canvas.getContext('2d');
		ctx.fillStyle = '#A35C02';
		ctx.fillRect(0, 0, canvas.width, canvas.height);

		let pathfields = {};
		if (path) {
			for (let field of path)
				pathfields[field.y * this.width + field.x] = true;
		}

		for (let y = 0; y < this.height; y++) {
			for (let x = 0; x < this.width; x++) {
				let i = y * this.width + x;
				if (pathfields[i]) {
					ctx.fillStyle = '#500';
					ctx.fillRect((1 + x) * DRAWSIZE, (1 + y) * DRAWSIZE, DRAWSIZE, DRAWSIZE);
				} else if (this.data[i] === CAVE_FREE) {
					ctx.fillStyle = '#E0AC4D';
					ctx.fillRect((1 + x) * DRAWSIZE, (1 + y) * DRAWSIZE, DRAWSIZE, DRAWSIZE);
				} else if (this.data[i] === CAVE_START) {
					ctx.fillStyle = '#ff0000';
					ctx.fillRect((1 + x) * DRAWSIZE, (1 + y) * DRAWSIZE, DRAWSIZE, DRAWSIZE);
				}
			}
		}
	}

	coordToCanvas(coord: number) {
		return Math.round((coord + 1.5) * DRAWSIZE);
	}
};
