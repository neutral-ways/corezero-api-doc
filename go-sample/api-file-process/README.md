# Upload-and-process

Uploads a CSV file into corezero and start the processing 

## Build

Run `make` to compile for linux and win

check `bin/` folder for binaries

```bash
make
```

## Config 

Update the `config.json` with the API-KEY and the API parameters 

## Run

Just pass a filename: 

```bash
./api-file-process test.csv
```

and the output:

```bash
Upload-and-process / CoreZero (c) 2022 🧟
API file upload started
 - API-KEY is: SVIjZXllML1hgEg0pVOQ3pikltoQEvlu
 - API host : 3.21.12.64

step 1: create upload request
 - attachment_id: 74cf11a4-aac6-426b-a48d-cdec4ac0e7c9
 - entity_id: 862c5e89-df30-48c2-95dc-3abb02d2595f
step 2: upload file
 - upload done. waiting 5 secs.....
step 3: process file
 - job created succesfully
 - job id: f34e0e49-6a3c-4c5f-822b-a61ebaefdecc
```



