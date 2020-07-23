import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import ParsedCave from '../shared/ParsedCave';
import { Treasure, VisitedField } from '../shared/Cave';

@Component({
	selector: 'app-cave-display',
	templateUrl: './cave-display.component.html',
	styleUrls: ['./cave-display.component.less']
})
export class CaveDisplayComponent implements OnInit {

	@Input() treasures: Treasure[] = [];

	_cave: ParsedCave = null;
	_path: VisitedField[] = [];

	@Input()
	set cave(c: ParsedCave) {
		this._cave = c;
		if (this._cave && this.cavecanvas && this.cavecanvas.nativeElement) {
			console.log(this.cavecanvas);
			this._cave.drawToCanvas(this.cavecanvas.nativeElement, this._path);
		} else {
			console.error('Could not draw to canvas!', this._cave, this.cavecanvas);
		}
	}

	@Input()
	set path(p: VisitedField[]) {
		this._path = p;
		this._cave.drawToCanvas(this.cavecanvas.nativeElement, this._path);
	}

	get cave() {
		return this._cave;
	}


	@ViewChild('cavecanvas', {static: true}) cavecanvas: ElementRef;

	@ViewChild('cavecontainer', {static: false}) cavecontainer: ElementRef;


	constructor() {
	}

	ngOnInit() {
	}

	focus(i: number) {
		console.log(this.cavecontainer.nativeElement.getElementsByClassName('marker-treasure')[i]);
		this.cavecontainer.nativeElement.getElementsByClassName('marker-treasure')[i].scrollIntoView(false);
	}

}
