import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { Cave } from '../shared/Cave';
import { BackendService } from '../shared/backend.service';
import { ActivatedRoute } from '@angular/router';
import ParsedCave from '../shared/ParsedCave';
import { BsModalRef, BsModalService } from 'ngx-bootstrap';
import { FormGroup, NgForm } from '@angular/forms';
import { MessageService } from '../shared/messages.service';
import { CaveDisplayComponent } from '../cave-display/cave-display.component';


@Component({
	selector: 'app-showcave',
	templateUrl: './showcave.component.html',
	styleUrls: ['./showcave.component.less']
})
export class ShowCaveComponent implements OnInit {

	cave: Cave = null;
	parsedCave: ParsedCave = null;

	treasureName: string;
	hideTreasureRef: BsModalRef;
	formHideTreasure: NgForm;

	@ViewChild('templateHideTreasure', {static: false}) templateHideTreasure: ElementRef;
	@ViewChild('cavedisplay', {static: false}) cavedisplay: CaveDisplayComponent;

	constructor(public backend: BackendService,
				private messages: MessageService,
				private modalService: BsModalService,
				private route: ActivatedRoute) {
	}

	ngOnInit() {
		this.getCave();
	}

	getCave(): void {
		const id = this.route.snapshot.paramMap.get('id');
		this.backend.getCave(id).subscribe(cave => {
			this.cave = cave;
			this.parsedCave = null;
			this.backend.getCaveTemplate(cave.template_id).subscribe(tmpl => this.parsedCave = new ParsedCave(tmpl));
		});
	}

	showHideTreasureModal() {
		this.hideTreasureRef = this.modalService.show(this.templateHideTreasure);
	}

	hideTreasure(form: NgForm) {
		if (!form.valid) {
			return false;
		}
		let name = this.treasureName;
		this.backend.hideTreasure(this.cave.id, this.treasureName).subscribe(cave => {
			this.cave = cave;
			this.messages.success('Treasure ' + name + ' has been hidden');
		});
		return true;
	}

}
