import {Injectable} from '@angular/core';
import Message from "./Message";

@Injectable()
export class MessageService {

	messages: Message[] = [];

	add(message: Message) {
		this.messages.push(message);
	}

	clear() {
		this.messages = [];
	}

	success(message: string) {
		this.add({type: 'success', text: message, timeout: 5000});
	}

	error(message: string) {
		this.add({type: 'danger', text: message, timeout: null});
	}

}
