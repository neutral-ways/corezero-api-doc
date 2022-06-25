# Corezero api integration 

## Abstract 

This documents defines how to interact with `corezero api` and generate transactions.

A transaction defines an operation to be quantified and its main goal is to generate reduction units. 

Each transaction contains a list of items and each one is quanitified -if its possible- 
based on its product type and other parameters defined by the platform. 


```json
{
    "account_id": "1bd46ad4-0bf3-4374-be98-06afb4e63d08",
    "transaction": {
        "reference": "ref-t50",
        "opened_at": "2022-01-02T14:10:10Z",
        "items": [
            {
                "reference": "popcorn-smallbag",
                "name": "salty popcorn",
                "quantity": 380,
                "product_line": "food",
                "product_value": 500,
                "product_unit": "mls",
                "factor_unit": "grs",
                "factor_value": 400,
                "source": "popcorn-factory"
            }
        ]
    }
}
```




## How it works 

- the API will try to quantify the transactions provided on a base effort basis
- in cases where sku/reference is not on the system the api will try to create the product
- in cases where the product_line is not on the system the api will try to create the product_line
- in cases where is not possible to categorize a product by product_line or category the operation status will be `pending` 



## Authentication

To use the transaction API an `API KEY` is needed.  

The API KEY is a unique key provided by corezero for each account. 

Any request to the API must include the key on the HTTP header `X-API-KEY` 

```bash
curl --location --request POST 'http://localhost:8080/api/v1/client/transaction' \
--header 'X-API-KEY: zaCELgL.0imfnc8mVLWwsAawjYr4RxAf50DDqtlx'
```

## Resources 

JSON objects used to compose a transaction request 
### Transaction 
A transaction is composed by a transaction header and a list of items


<!-- { "blockType": "resource",
"@type": "transaction",
"optionalProperties": [] } -->
```json
{
  "reference": "string",
  "opened_at": "datetime",
  "items": [ {"@type": "transaction_item"} ]
}
```


### Properties

| Property | Type | Description | Mandatory |
|:---------|:-----|:------------|:----------|
| **reference**   | String | Unique identifier for the TX in the lot. | yes |
| **opened_at** | String | The date the operation took place in ISO-8601. | yes |
| **project_id**| uuid | The id of the project. If not provided the API will look for any open project.  | no | 
| **item** | Collection([transaction_item](#transaction-item)) | A collection of transaction itmes. | yes |


Some validation rules: 
- The project_id must be an `uuid` v4 value
- If `project_id` is not provided and more than one projects are open on the account the api will throw an error
- `opened_at` can be a date in the past but cannot be a date in the future
- At least one transaction item is required on items


### Transaction Item

<!-- { "blockType": "resource",
"@type": "transaction_item",
"optionalProperties": [] } -->
```json
{
  "quantity": "string",
  "sku": "string",
  "reference": "string",
  "name": "datetime",
  "unit": "string",
  "product_value": "string",
  "source": "string",
  "product_line": "string",
  "factor_unit": "string",
  "factor_value": "string",
  "product_unit": "string" 
}
```

| Property | Type | Description | Mandatory |
|:---------|:-----|:------------|:----------|
| **quantity**   | float | Amount of units to be quantified. cannot be 0 | yes |
| **sku** | String | Identifier of the product in the inventory | no |
| **reference**| uuid | Customer external identifier for the product | no | 
| **name** | String | The name of the product | no |
| **unit** | string |  |
| **product_value** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **source** | String | The source or seller of the product | no |
| **product_line** | String | the name of the product line for the product | no |
| **factor_unit** | String | The unit type for the factor. see ([catalog](#unit-types-catalog)) | no |
| **factor_value** | String | Factor value for quantification | no |
| **product_unit** | String | Product unit type. see ([catalog](#unit-types-catalog)). | no |


Some validation rules: 
- Quantity must be greather than 0
- SKU and Reference are not mandatory but you need at least one of them defined
- Name is only needed if the SKU/Rereference is not valid (product doesnt exsist)
- Attribute `opened_at` can be a date in the past but cannot be a date in the future
- At least one transaction item is required on items



# Create transaction API

Post a transaction into corezero middleware and try to process it. 

This endpoint exectues the task synchronically so you will get the response of each item immediately after the call. 

### Prerequisites

One of the following scopes are required to execute this request:

* To call this API you must have an API key (check authorization section)
* Each transaction has a reference id which is unique in the scope of the lot


**URL** : `/api/v1/client/transaction`

**METHOD** : `POST`

**Auth required** : `Yes / API-KEY`


### HTTP Request

```bash
curl --location --request POST "https://api-dev.corezz.net/api/v1/client/transaction" \
--header "X-API-KEY: a7da7.123ads7a8d7823" \
--header "Content-Type: application/json" \
--data-raw '{
    "reference": "ref-t50",
    "opened_at": "2022-01-02T14:10:10Z",
    "items": [
        {
            "reference": "invoice-6678",
            "name": "popcorn-small",
            "quantity": 4,
            "product_line": "food-mrkt",
            "product_value": 500,
            "product_unit": "mls",
            "factor_unit": "grs",
            "factor_value": 206,
            "source": "popcorn-factory"
        }
    ]

}'
```

### Request parameters

No parameters on the query string 

### Optional request headers

| Name | Value |
|:-----|:------|

### Request body

Do not supply a request body with this method.

In the request URL, provide the following query parameters with values.

| Property | Type | Description |
|:---------|:-----|:------------|
| **reference**   | String | Unique identifier for the TX in the lot. This field is mandatory |
| **opened_at** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **project_id**| uuid | The id of the project. If not provided the API will look for any open project.  | 
| **item** | Collection([transaction_item](#transaction-item)) | A collection of transaction itmes. |



## Success Responses

**Condition** : Data provided is valid and User is Authenticated.

**Code** : `200 OK`

**Content example** : Response will return a list of the items in the transaction. 
And the details about each processed item indicating if the operation was quantified or not. 


```json
{
    "data": {
        "account_id": "f7e56912-2dfe-4ec6-ab78-917ae5537800",
        "reference": "",
        "operation_status": "pending",
        "transaction": {
            "id": "a2f70cc1-0e65-4de7-afc6-0a0fb63dfc55",
            "created_at": "2022-05-17T17:35:13.8431658-03:00",
            "lot_id": "92d0cadf-443e-4c52-9df2-e0fe12df8034",
            "transaction_type": "in-api",
            "reference": "ref-t50",
            "opened_at": "2022-01-02T14:10:10Z",
            "status": "pending",
            "items": [
                {
                    "id": "fa62064c-4d6a-4ca1-bec3-5f88a130b36c",
                    "product_id": "c844a06e-943a-454e-9b74-d35b73e9a2b2",
                    "quantity": 4,
                    "factor": 0,
                    "quantification": 0,
                    "product_category_id": null,
                    "factor_unit_id": 1,
                    "factor_quantity": 206,
                    "quantification_by": "none"
                },
                {
                    "id": "8006e10c-aeff-4ef1-b0f6-0f6d7fc4a7c2",
                    "product_id": "23736e39-5f96-4107-b4c2-0687abfb350b",
                    "quantity": 4,
                    "factor": 1.8,
                    "quantification": 1483.2,
                    "product_category_id": null,
                    "factor_unit_id": null,
                    "factor_quantity": 206,
                    "quantification_by": "product_line"
                }
            ]
        }
    }
}
```

**Remarks**

- A transaction response may have two different status, if all of the items are quantified then the overall
status of the transaction will be `quantified`. If at least one item is pending, then the status is `pending` 

- `quantification_by` is how the middleware was able to quantify the item. By product category, by product line or none if it wasnt able to quantify.

- If an item is quantified on the response it will have all the attirbutes filled (`quantification`, `factor` and `quantification_by`)


## Error Response

**Condition** : If provided data is invalid or there is something stopping the transaction to be processed then error is thrown.

**Code** : `400 BAD REQUEST`

**Content example** :

Parameter validation error example: 

```json
{
    "errors": [
        {
            "field": "reference",
            "reason": "required",
            "error": "Key: 'ApiTransaction.reference' Error:Field validation for 'reference' failed on the 'required' tag"
        }
    ]
}
```


Process error example:

```json
{
    "error": 
        "reference is required",
}
```


### List of possible errors

| Message | Description |
|:---------|:------------|
| reference is required | attribute reference is missing |
| opened_at is required | attribute opened_at is missing or wrong |
| account is invalid | the requested account is not valid or doesnt exsist |
| no open project found | |
| more than one open project found |  |
| more than one open lot found | |
| no open lot found | |
| name is required for creating a new product | |
| tx reference alredy exists| the reference of tx you are trying to create already exsist |



## Unit types catalog


| Unit | Description |
|:---------|:------------|
| **ml**   | Mililiters |
| **lt**   | Liters |
| **km2**   | Square kilometers |
| **cm**   | Centimeters |
| **mm**   | Milimeters |
| **km**   | kilometers |
| **grs**   | grams |
| **m3**   | cubic meters |
| **u**   | units |
| **kg**   | Kilograms |
| **t**   | Tones |
| **fg**   | frigoria |



# File upload API 

Using this API endpoints client can upload a CSV file with transaction in a predefined format 

and order the middleware to process the file asynchronically. 


## How it works?

The process has 3 steps: 

1. First you create an upload request where you need to pass the filename and the content-type (currently the only supported is `text/csv`)

2. You upload the file to the URL obtained from the step 1 call. You must include the headers returned from step 1. 

3. Once the file is uploaded you tell the API that you want to process the file.

![arch](./seq-diagram-upload.png)




## Step 1: create upload request


On this step you will tell the API you want to upload a file


**URL** : `/api/v1/client/file`

**METHOD** : `POST`

**Auth required** : `Yes / API-KEY`


### HTTP Request

```bash
curl --location --request POST 'http://localhost:8080/api/v1/client/file' \
--header 'X-API-KEY: SVIjZXllML1hgEg0pVOQ3pikltoQEvlu' \
--header 'Content-Type: application/json' \
--data-raw '{
    "content_type": "text/csv",
    "filename": "sugo-co-january-2021.csv"
}'
```

### Request body

| Property | Type | Description |
|:---------|:-----|:------------|
| **content_type**   | String | Content type of the file. This field is mandatory. The only supported is `text/csv` |
| **filename** | String | The name of the file you are uploading. This field is mandatory |


### HTTP Response


```json
{
    "data": {
        "id": "51c1eb60-fae1-477a-97e9-b914400cc841",
        "entity": "tx-processor",
        "entity_id": "862c5e89-df30-48c2-95dc-3abb02d2595f",
        "filename": "sugo-co-january-2021.csv",
        "pre_signed_url": "https://files-corezero.s3.us-east-2.amazonaws.com/tx-processor/862c5e89-df30-48c2-95dc-3abb02d2595f/51c1eb60-fae1-477a-97e9-b914400cc841.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAVEHI6TFWTYQT42TM%2F20220625%2Fus-east-2%2Fs3%2Faws4_request&X-Amz-Date=20220625T142413Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host%3Bx-amz-meta-content-type%3Bx-amz-meta-entity%3Bx-amz-meta-entity-id%3Bx-amz-meta-filename%3Bx-amz-meta-public%3Bx-amz-meta-uploader&x-id=PutObject&X-Amz-Signature=7e58861f2aec3171e0d8155be989f850b07bca3f77610d54c416e4648a9439e8",
        "headers": {
            "x-amz-meta-content-type": "text/csv",
            "x-amz-meta-entity": "tx-processor",
            "x-amz-meta-entity-id": "862c5e89-df30-48c2-95dc-3abb02d2595f",
            "x-amz-meta-filename": "sugo-co-january-2021.csv",
            "x-amz-meta-public": "false",
            "x-amz-meta-uploader": "api-key"
        }
    }
}
```

| Property | Type | Description |
|:---------|:-----|:------------|
| **id**   | String | Is the id of the file once uploaded (If upload is succesfull)  |
| **entity**   | String | Internal reference to group the uploaded files |
| **entity_id** | String | Is the account id related to the file. Its taken from the api-key |
| **filename** | String | The name of the file you are uploading. |
| **pre_signed_url** | String | The URL where  you have to Upload the file on the next step. |
| **headers** | String | List of headers you need to include on the upload request on the next step. |




## Step 2: Upload the file to the storage


On this step you will upload the file to the storage. 
The URL you must is the one returned on the `step 1`. 

**URL** : `{pre-signed-url}`

**METHOD** : `PUT`

**Auth required** : `NO`


### HTTP Request

```bash
curl --location --request PUT 'https://files-corezero.s3.us-east-2.amazonaws.com/tx-processor/e8b93001-d531-4f98-a2c4-8bee079f0ff2/8939dbd0-bc60-4f22-acc7-69a4539c5935.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAVEHI6TFWTYQT42TM%2F20220625%2Fus-east-2%2Fs3%2Faws4_request&X-Amz-Date=20220625T143424Z&X-Amz-Expires=900&X-Amz-SignedHeaders=host%3Bx-amz-meta-content-type%3Bx-amz-meta-entity%3Bx-amz-meta-entity-id%3Bx-amz-meta-filename%3Bx-amz-meta-public%3Bx-amz-meta-uploader&x-id=PutObject&X-Amz-Signature=2b83c92663c9ba81f7adbcd4427fe09519c995940ee128cb2cbb082aa0424c4a' \
--header 'x-amz-meta-entity: tx-processor' \
--header 'x-amz-meta-entity-id: e8b93001-d531-4f98-a2c4-8bee079f0ff2' \
--header 'x-amz-meta-filename: con-p-malo.csv' \
--header 'x-amz-meta-uploader: gipshm+test@gmail.com' \
--header 'x-amz-meta-content-type: text/csv' \
--header 'x-amz-meta-public: false' \
--header 'Content-Type: text/csv' \
--data-binary '@/C:/Users/gipsh/Downloads/sugo-co-january-2021.csv'
```



### HTTP Response

The response is only `200` if everything went well.

If something went wrong you will have an error from the server exaplaining the problem. 

Most common issue is a signature error meaning something is worng with headers. 

```xml
<?xml version="1.0" encoding="UTF-8"?>
<Error>
    <Code>SignatureDoesNotMatch</Code>
    <Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
    <AWSAccessKeyId>AKIAVEHI6TFWTYQT42TM</AWSAccessKeyId>
</Error>
```



## Step 3: process the file


On this step you will instruct the middleware to process the uploaded file asynchronically


**URL** : `/api/v1/client/process`

**METHOD** : `POST`

**Auth required** : `Yes / API-KEY`


### HTTP Request

```bash
curl --location --request POST 'http://3.21.12.64:8080/api/v1/client/process' \
--header 'X-API-KEY: SVIjZXllML1hgEg0pVOQ3pikltoQEvlu' \
--header 'Content-Type: application/json' \
--data-raw '{
    "attachment_id": "ae3a9167-99f3-45e4-bafa-94fef0d3d1aa"
}'
```

### Request body

| Property | Type | Description |
|:---------|:-----|:------------|
| **attachment_id**   | String | Is the id of the uploaded file. The filed `id` returned on step 1. This field is mandatory. |


### HTTP Response

If attachment_id is correct it will respond with:

```json
{
    "data": {
        "id": "86bdd8a4-9de6-4b5f-b399-0620723c60b7",
        "created_at": "2022-06-25T15:06:38.958072414Z",
        "updated_at": "2022-06-25T15:06:38.958072414Z",
        "deleted_at": null,
        "created_by": "gipshm@gmail.com",
        "updated_by": "",
        "account_id": null,
        "project_id": null,
        "lot_id": null,
        "attachment_id": "d6b6930b-af0e-4bd7-9edb-7f649df39642",
        "status": "created",
        "error_message": null,
        "ended_at": null,
        "operation": "process"
    }
}
```

| Property | Type | Description |
|:---------|:-----|:------------|
| **id**   | String | The id of the running job  |
| **created_by**   | String | The user who created the job or `api-key` |
| **entity_id** | String | Is the account id related to the file. Its taken from the api-key |
| **attachment_id** | String | The id of the file used in the process. |
| **status** | String | The status of the operation. Its going to be `created`. |
| **operation** | String | The type of the operation. In this case is going to be `process`. |

