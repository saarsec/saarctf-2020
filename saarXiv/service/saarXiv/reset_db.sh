#!/bin/bash
set -eux
rm -f app.db
rm -f Migrations/*
PATH=$PATH:/home/johannes/.dotnet/tools dotnet ef migrations add InitialCreate
PATH=$PATH:/home/johannes/.dotnet/tools dotnet ef database update

