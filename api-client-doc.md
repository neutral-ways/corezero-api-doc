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
- in cases where sku/refenrece is not on the system the api will try to create the product
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

## Model 

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
- `opened_at` can be a date in the pasto but cannot be a date in the future
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
| **unit** | string | A collection of author resources that represent the authors of the book. |
| **product_value** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **source** | String | The source or seller of the product | no |
| **product_line** | String | the name of the product line for the product | no |
| **factor_unit** | String | The unit type for the factor. see catalog | no |
| **factor_value** | String | Factor value for quantification | no |
| **product_unit** | String | Product unit size. see catalog. | no |


Some validation rules: 
- Quantity must be greather than 0
- SKU and Reference are not mandatory but you need at least one of them defined
- Name is only needed if the SKU/Rereference is not valid (product doesnt exsist)
- 
- `opened_at` can be a date in the pasto but cannot be a date in the future
- At least one transaction item is required on items


## Transaction API 



# Create transaction

Creates a transaction in corezero API to be processed. 
This method exectues synchronically so you will get the response of each item immediately after the call. 

### Prerequisites

One of the following scopes are required to execute this request:

* To call this API you must have an API key (check authorization section)
* Each transaction has a reference id which is unique in the scope of the lot


** URL ** : `/api/v1/client/transaction`

** METHOD ** : `POST`

** Auth required ** : YES / API-KEY


### HTTP Request

```json
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
| **u**   |  |
| **kg**   | Kilograms |
| **t**   | Tones |
| **fg**   | frigoria |


