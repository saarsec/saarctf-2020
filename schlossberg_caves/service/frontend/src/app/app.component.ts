import { Component } from '@angular/core';

import { ApiService } from './shared';

@Component({
	selector: 'my-app', // <my-app></my-app>
	templateUrl: './app.component.html',
	styleUrls: ['app.component.less'],
})
export class AppComponent {
	title: string;

	constructor(private api: ApiService) {
		this.title = this.api.title;
	}
}
