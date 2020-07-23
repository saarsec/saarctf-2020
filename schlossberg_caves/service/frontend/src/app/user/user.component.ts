import { Component, OnInit } from '@angular/core';
import { BackendService } from '../shared/backend.service';

@Component({
	selector: 'app-user',
	templateUrl: './user.component.html',
	styleUrls: []
})
export class UserComponent implements OnInit {

	user = null;
	username: string = '';
	password: string = '';

	constructor(public backend: BackendService) {
	}

	ngOnInit() {
		this.backend.getCurrentUser().subscribe(user => this.user = user);
	}

	login(username, password) {
		this.backend.login(username, password);
	}

	register(username, password) {
		this.backend.register(username, password);
	}

	logout() {
		this.backend.logout();
	}

}
