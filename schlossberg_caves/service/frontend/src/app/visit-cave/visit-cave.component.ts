import { Component, OnInit, ViewChild } from '@angular/core';
import { BackendService } from '../shared/backend.service';
import { BsModalService, TabsetComponent } from 'ngx-bootstrap';
import ParsedCave from '../shared/ParsedCave';
import { MessageService } from '../shared/messages.service';
import { ActivatedRoute } from '@angular/router';
import { Cave, Treasure, VisitedField } from '../shared/Cave';
import { setTime } from 'ngx-bootstrap/timepicker/timepicker.utils';


const defaultCode = `
eija main() gebbtserick int: {
	runner;
	riwwer;
	ferdisch;
	serick 0;
}
`;


@Component({
	selector: 'app-visit-cave',
	templateUrl: './visit-cave.component.html',
	styleUrls: ['./visit-cave.component.less']
})
export class VisitCaveComponent implements OnInit {

	@ViewChild('filetabs', {static: false}) filetabs: TabsetComponent;
	files: string[] = ['entry.sl'];
	codefiles: { [key: string]: string } = {'entry.sl': defaultCode};

	codeConfig = {indentUnit: 4, tabSize: 4, indentWithTabs: true};

	output = '';
	success = false;
	visitedPath: { path: VisitedField[], treasures: Treasure[] } = {path: [], treasures: []};

	cave: Cave = null;
	parsedCave: ParsedCave = null;


	constructor(private backend: BackendService,
				private messages: MessageService,
				private route: ActivatedRoute) {
	}

	ngOnInit() {
		this.getCave();
		const handler = () => {
			if (this.filetabs) {
				this.reorderTabs();
				this.filetabs.tabs[0].active = true;
				console.log(this.filetabs.tabs[0]);
				this.initTab(0);
			} else {
				setTimeout(handler, 50);
			}
		};
		setTimeout(handler, 50);
	}

	getCave(): void {
		const id = this.route.snapshot.paramMap.get('id');
		this.backend.getCave(id).subscribe(cave => {
			this.cave = cave;
			this.parsedCave = null;
			this.backend.getCaveTemplate(cave.template_id).subscribe(tmpl => this.parsedCave = new ParsedCave(tmpl));
		});
	}

	reorderTabs() {
		if (!this.filetabs) return;
		this.filetabs.tabs.sort((a, b) => {
			if (a.customClass === 'addFileTab') {
				return 1;
			} else if (b.customClass === 'addFileTab') {
				return -1;
			}
			return a.heading.localeCompare(b.heading);
		})
	}

	addFile() {
		let filename = prompt('File name', '');
		if (filename) {
			if (this.files.indexOf(filename) < 0) {
				this.files.push(filename);
			}
			if (!this.codefiles[filename]) {
				this.codefiles[filename] = '// ' + filename + '\n\n';
			}
		}
		setTimeout(() => {
			this.reorderTabs();
			this.filetabs.tabs[this.filetabs.tabs.length - 2].active = true;
			this.initTab(this.filetabs.tabs.length - 2);
		}, 1);
	}


	removeFile(filename) {
		console.log('REMOVE', filename);
		this.files.splice(this.files.indexOf(filename), 1);
		delete this.codefiles[filename];
		setTimeout(() => this.reorderTabs(), 1);
	}


	private initTab(index) {
		const x = this.filetabs.tabs[index].elementRef.nativeElement.querySelector('.CodeMirror-scroll');
		setTimeout(() => {
			x.dispatchEvent(new MouseEvent('mousedown', {bubbles: true, cancelable: true, view: window}));
		}, 10);
	}


	runCode() {
		this.backend.executeCode(this.cave.id, this.codefiles).subscribe(result => {
			this.output = result;
			this.parseOutput();
		});
	}

	parseOutput() {
		let p = this.output.indexOf('VISITED PATH: ');
		this.success = p > 0;
		if (this.success) {
			let msg = 'Visit finished - did you enjoy it?';
			this.visitedPath = JSON.parse(this.output.substr(p + 14));
			this.output = this.output.substr(0, p);

			if (this.visitedPath.treasures.length === 0) {
				msg += ' Unfortunately we didn\'t see any treasures.';
			} else {
				msg += ' You found ' + this.visitedPath.treasures.length + ' treasures:';
				for (let t of this.visitedPath.treasures)
					msg += ' - ' + t.name + ' (' + t.x + ',' + t.y + ')';
			}

			this.messages.success(msg);
		} else {
			this.messages.error('Your guide got lost in the cave. We\'re sorry.');
		}
	}
}
