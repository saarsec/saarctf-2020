export class Treasure {
	name: string;
	x: number;
	y: number;
}

export class VisitedField {
	x: number;
	y: number;
}

export class Cave {
	id: string;
	name: string;
	owner: string;
	template_id: number;
	created: number;
	treasure_count: number;
	treasures?: Treasure[];
}
