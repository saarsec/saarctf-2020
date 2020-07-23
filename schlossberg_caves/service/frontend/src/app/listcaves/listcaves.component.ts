import { Component, OnInit } from '@angular/core';
import { BackendService } from '../shared/backend.service';
import { Cave } from '../shared/Cave';

@Component({
	selector: 'app-listcaves',
	templateUrl: './listcaves.component.html',
	styleUrls: []
})
export class ListCavesComponent implements OnInit {

	caves: Cave[];

	constructor(private backend: BackendService) {
	}

	ngOnInit() {
		this.backend.getCaves().subscribe(caves => this.caves = caves);
	}

}
