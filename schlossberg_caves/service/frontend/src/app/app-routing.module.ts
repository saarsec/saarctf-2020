import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from './home/home.component';
import { RentCaveComponent } from './rentcave/rentcave.component';
import { ShowCaveComponent } from './showcave/showcave.component';
import { ListCavesComponent } from './listcaves/listcaves.component';
import { VisitCaveComponent } from './visit-cave/visit-cave.component';

const routes: Routes = [
	{path: '', component: HomeComponent},
	{path: 'rent-a-cave', component: RentCaveComponent},
	{path: 'cave/:id', component: ShowCaveComponent},
	{path: 'cave/:id/visit', component: VisitCaveComponent},
	{path: 'caves', component: ListCavesComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
