import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { MessageService } from './shared/messages.service';
import { BackendService } from './shared/backend.service';
import { ApiService } from './shared/api.service';

import { HomeComponent } from './home/home.component';
import { RentCaveComponent } from './rentcave/rentcave.component';
import { ShowCaveComponent } from './showcave/showcave.component';
import { ListCavesComponent } from './listcaves/listcaves.component';
import { VisitCaveComponent } from './visit-cave/visit-cave.component';
import { CaveDisplayComponent } from './cave-display/cave-display.component';
import { MessagesComponent } from './messages/messages.component';
import { UserComponent } from './user/user.component';


import { AlertModule, ModalModule, TabsModule } from 'ngx-bootstrap';
import { DragScrollModule } from 'ngx-drag-scroll';
import { CodemirrorModule } from 'ng2-codemirror';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    RentCaveComponent,
    ShowCaveComponent,
    ListCavesComponent,
    VisitCaveComponent,
    CaveDisplayComponent,
    MessagesComponent,
    UserComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule,
    DragScrollModule,
	ModalModule.forRoot(),
	AlertModule.forRoot(),
	TabsModule.forRoot(),
	CodemirrorModule,
  ],
  providers: [
  	MessageService,
    BackendService,
    ApiService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
