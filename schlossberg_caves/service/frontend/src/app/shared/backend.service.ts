import { Injectable } from '@angular/core';
import { CaveTemplate } from './CaveTemplate';
import { MessageService } from './messages.service';
import { Observable, of } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { catchError, map } from 'rxjs/operators';
import { Cave } from './Cave';


interface UserResponse {
	username: string;
}


@Injectable()
export class BackendService {

	private backendUrl = 'api/';

	constructor(private http: HttpClient, private messageService: MessageService) {
	}

	private getDefaultErrorHandler<T>(operation = 'operation', result?: T) {
		return (error: any): Observable<T> => {
			console.error(error);
			this.messageService.error(`${operation} failed: ${error.message}`);
			// Let the app keep running by returning an empty result.
			return of(result as T);
		};
	}

	private caveTemplates: Observable<CaveTemplate[]> = null;

	getCaveTemplates(): Observable<CaveTemplate[]> {
		if (this.caveTemplates === null) {
			this.caveTemplates = this.http.get<CaveTemplate[]>(this.backendUrl + 'templates/list').pipe(
				catchError(this.getDefaultErrorHandler('getCaveTemplates', []))
			);
		}
		return this.caveTemplates;
	}

	getCaveTemplate(id: number) {
		return this.http.get(this.backendUrl + 'templates/' + id,
			{responseType: 'arraybuffer'}).pipe(
			catchError(this.getDefaultErrorHandler('getCaveTemplate',))
		);
	}


	public username: string = null;

	login(username: string, password: string) {
		this.http.post<UserResponse>(this.backendUrl + 'users/login', {username: username, password: password})
			.subscribe(user => {
					this.username = user.username;
					this.messageService.success('Welcome, ' + user.username);
				},
				error => {
					console.error(error);
					this.messageService.error('Invalid username or password');
				});
	}

	register(username: string, password: string) {
		this.http.post<UserResponse>(this.backendUrl + 'users/register', {username: username, password: password})
			.subscribe(user => {
					this.username = user.username;
					this.messageService.success('Welcome, ' + user.username);
				},
				error => {
					console.error(error);
					this.messageService.error(error.error);
				});
	}

	getCurrentUser(): Observable<string> {
		let response = this.http.get<UserResponse>(this.backendUrl + 'users/current')
			.pipe(map(user => user.username))
			.pipe(catchError(this.getDefaultErrorHandler('getCurrentUser', null)));
		response.subscribe(user => {
			this.username = user;
		});
		return response;
	}

	logout() {
		this.http.post(this.backendUrl + 'users/logout', {})
			.pipe(catchError(this.getDefaultErrorHandler('logout')))
			.subscribe(_ => this.username = null);
	}


	rentCave(templateId: number, name: string): Observable<Cave> {
		return this.http.post<Cave>(this.backendUrl + 'caves/rent', {
			template: templateId,
			name: name
		}).pipe(
			catchError(this.getDefaultErrorHandler('rentCave', null))
		);
	}

	getCaves(): Observable<Cave[]> {
		return this.http.get<Cave[]>(this.backendUrl + 'caves/list').pipe(
			catchError(this.getDefaultErrorHandler('getCaves', []))
		);
	}

	getCave(id: string): Observable<Cave> {
		return this.http.get<Cave>(this.backendUrl + 'caves/' + id).pipe(
			catchError(this.getDefaultErrorHandler('getCave', null))
		);
	}

	hideTreasure(caveId: string, name: string): Observable<Cave> {
		return this.http.post<Cave>(this.backendUrl + 'caves/hide-treasures', {
			cave_id: caveId,
			names: [name]
		}).pipe(
			catchError(this.getDefaultErrorHandler('hideTreasure', null))
		);
	}

	executeCode(caveId: string, files: { [key: string]: string }): Observable<string> {
		return this.http.post(this.backendUrl + 'visit', {cave_id: caveId, files: files}, {responseType: 'text'})
			.pipe(catchError(this.getDefaultErrorHandler('executeCode', '')));
	}

}
