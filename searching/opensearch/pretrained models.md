## Supported pretrained models
 - Text embedding models are sourced from Hugging Face
 - Sparse encoding models are trained by OpenSearch
### Sentence transformers
### Sparse encoding models
### Cross-encoder models

## Prerequisites
- Setup enable 

curl -k -XPUT "https://localhost:9200/_cluster/settings" -H 'Content-Type: application/json' -d'
{
  "persistent": {
    "plugins": {
      "ml_commons": {
        "only_run_on_ml_node": "false",
        "model_access_control_enabled": "true",
        "native_memory_threshold": "99"
      }
    }
  }
}' -u admin:admin

+ Result: 
{"acknowledged":true,"persistent":{"plugins":{"ml_commons":{"only_run_on_ml_node":"false","model_access_control_enabled":"true","native_memory_threshold":"99"}}},"transient":{}}

## Step 1: Register a model group
- curl -k -XPOST "https://localhost:9200/_plugins/_ml/model_groups/_register" -H 'Content-Type: application/json' -d'
{
  "name": "local_model_group",
  "description": "A model group for local models"
}' -u admin:admin

+ Result: 
{"model_group_id":"xIkZEJABlLB-9JtUgosr","status":"CREATED"}

## Step 2: Register a local OpenSearch-provided model
- curl -k -XPOST "https://localhost:9200/_plugins/_ml/models/_register" -H 'Content-Type: application/json' -d'
{
  "name": "huggingface/sentence-transformers/msmarco-distilbert-base-tas-b",
  "version": "1.0.2",
  "model_group_id": "xIkZEJABlLB-9JtUgosr",
  "model_format": "TORCH_SCRIPT"
}' -u admin:admin 
+ Result: 
{"task_id":"xYkaEJABlLB-9JtUXYsE","status":"CREATED"}

- curl -k -XGET "https://localhost:9200/_plugins/_ml/tasks/xYkaEJABlLB-9JtUXYsE" -u admin:admin 

+ Result: 

{"model_id":"xokaEJABlLB-9JtUZovQ","task_type":"REGISTER_MODEL","function_name":"TEXT_EMBEDDING","state":"COMPLETED","worker_node":["OHC8dRlPQRqq1q7-uTbYuA"],"create_time":1718257081491,"last_update_time":1718257121393,"is_async":true}% (base) 

## Step 3: Deploy the model
- curl -k -XPOST "https://localhost:9200/_plugins/_ml/models/xokaEJABlLB-9JtUZovQ/_deploy"  -u admin:admin 

+ Result: 
{"task_id":"yIkeEJABlLB-9JtU34uU","status":"CREATED"}% 


- curl -k -XGET "https://localhost:9200/_plugins/_ml/tasks/yIkeEJABlLB-9JtU34uU"  -u admin:admin 

+ result: 
{"model_id":"xokaEJABlLB-9JtUZovQ","task_type":"DEPLOY_MODEL","function_name":"TEXT_EMBEDDING","state":"COMPLETED","worker_node":["OHC8dRlPQRqq1q7-uTbYuA"],"create_time":1718257377168,"last_update_time":1718257404499,"is_async":true}

## Step 4 (Optional): Test the model
- For a text embedding model
  
curl -k -XPOST "https://localhost:9200/_plugins/_ml/_predict/text_embedding/xokaEJABlLB-9JtUZovQ" -H 'Content-Type: application/json' -d'
{
  "text_docs":[ "today is sunny"],
  "return_number": true,
  "target_response": ["sentence_embedding"]
}' -u admin:admin


- For a parse encoding model 
curl -XPOST "http://localhost:9200/_plugins/_ml/_predict/sparse_encoding/cleMb4kBJ1eYAeTMFFg4" -H 'Content-Type: application/json' -d'
{
  "text_docs":[ "today is sunny"]
}'

- For a question/answer model 
curl -XPOST "http://localhost:9200/_plugins/_ml/models/<model_id>/_predict" -H 'Content-Type: application/json' -d'
{
    "query_text": "today is sunny",
    "text_docs": [
        "how are you",
        "today is sunny",
        "today is july fifth",
        "it is winter"
    ]
}'
## Step 5: Use the model for search
    - Semetic search 
    - ....

