import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { BackendService } from '../shared/backend.service';
import { CaveTemplate } from '../shared/CaveTemplate';
import { BsModalRef, BsModalService } from 'ngx-bootstrap';
import ParsedCave from '../shared/ParsedCave';
import { Router } from '@angular/router';

@Component({
	selector: 'app-rentcave',
	templateUrl: './rentcave.component.html',
	styleUrls: []
})
export class RentCaveComponent implements OnInit {

	@ViewChild('templateCavePreview', {static: false}) templateCavePreview: ElementRef;

	public templates: CaveTemplate[];

	constructor(private backend: BackendService, private modalService: BsModalService, private router: Router) {
	}

	ngOnInit() {
		this.backend.getCaveTemplates().subscribe(templates => this.templates = templates);
	}


	cavePreviewRef: BsModalRef;
	currentCave: ParsedCave;
	currentCaveTemplate: CaveTemplate;

	previewCave(template: CaveTemplate) {
		if (!this.backend.username) {
			alert('Please log in first!');
			return;
		}
		this.currentCave = null;
		this.currentCaveTemplate = template;
		this.backend.getCaveTemplate(template.id).subscribe(cave => this.currentCave = new ParsedCave(cave));
		this.cavePreviewRef = this.modalService.show(this.templateCavePreview, {'class': 'large'});
	}


	rentCave(cavetmpl: CaveTemplate) {
		this.backend.rentCave(cavetmpl.id, cavetmpl.name).subscribe(cave => {
			if (cave) {
				this.router.navigate(['/cave/' + cave.id]);
				if (this.cavePreviewRef) {
					this.cavePreviewRef.hide();
				}
			}
		});
	}

}
