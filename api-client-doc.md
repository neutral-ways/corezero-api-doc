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

| Property | Type | Description |
|:---------|:-----|:------------|
| **reference**   | String | Unique identifier for the TX in the lot. This field is mandatory |
| **opened_at** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **project_id**| uuid | The id of the project. If not provided the API will look for any open project.  | 
| **item** | Collection([transaction_item](#transaction-item)) | A collection of transaction itmes. |


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
  "product_value": "string"
  "source": "string"
  "product_line": "string"
  "factor_unit": "string"
  "factor_value": "string"
  "product_unit": "string"
  
}
```

| Property | Type | Description |
|:---------|:-----|:------------|
| **quantity**   | float | Amount of units to be quantified. Is mandatory and cannot be 0 |
| **sku** | String | Identifier of the product in the inventory |
| **reference**| uuid | Customer external identifier for the product | 
| **name** | String | The name of the product |
| **unit** | Collection([transaction_item]) | A collection of author resources that represent the authors of the book. |
| **product_value** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **source** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **product_line** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **factor_unit** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **factor_value** | String | The date the operation took place in ISO-8601. This field is mandatory |
| **product_unit** | String | The date the operation took place in ISO-8601. This field is mandatory |



## Transaction API 


