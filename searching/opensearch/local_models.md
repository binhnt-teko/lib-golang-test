## Custom Local Models
### model support 
    - 2.6: local text embedding models.
    - 2.11:  local sparse encoding models
    - 2.12: local cross-encoder models
    - 2.13:  local question answering models.
### Preparing a model
    - provide a tokenizer JSON file within the model zip file
    - sparse encoding models, make sure your output format is {"output":<sparse_vector>} so that ML Commons can post-process the sparse vector
    - Model format:
      - portable format: TorchScript and ONNX formats
      - must calculate a SHA256 checksum for the model zip file: shasum -a 256 sentence-transformers_paraphrase-mpnet-base-v2-1.0.0-onnx.zip
    - Model size:
      - splits the model file into smaller chunks to be stored in a model index
### Prerequisites
    - Cluster settings:
  curl -XPUT "https://localhost:9200/_cluster/settings" -H 'Content-Type: application/json' -d'
{
  "persistent": {
    "plugins": {
      "ml_commons": {
        "allow_registering_model_via_url": "true",
        "only_run_on_ml_node": "false",
        "model_access_control_enabled": "true",
        "native_memory_threshold": "99"
      }
    }
  }
}'

### Step 1: Register a model group 
  curl -XPOST "http://localhost:9200/_plugins/_ml/model_groups/_register" -H 'Content-Type: application/json' -d'
{
  "name": "local_model_group",
  "description": "A model group for local models"
}'

### Step 2: Register a local model
curl -XPOST "http://localhost:9200/_plugins/_ml/models/_register" -H 'Content-Type: application/json' -d'
{
  "name": "huggingface/sentence-transformers/msmarco-distilbert-base-tas-b",
  "version": "1.0.1",
  "model_group_id": "wlcnb4kBJ1eYAeTMHlV6",
  "description": "This is a port of the DistilBert TAS-B Model to sentence-transformers model: It maps sentences & paragraphs to a 768 dimensional dense vector space and is optimized for the task of semantic search.",
  "model_task_type": "TEXT_EMBEDDING",
  "model_format": "TORCH_SCRIPT",
  "model_content_size_in_bytes": 266352827,
  "model_content_hash_value": "acdc81b652b83121f914c5912ae27c0fca8fabf270e6f191ace6979a19830413",
  "model_config": {
    "model_type": "distilbert",
    "embedding_dimension": 768,
    "framework_type": "sentence_transformers",
    "all_config": "{\"_name_or_path\":\"old_models/msmarco-distilbert-base-tas-b/0_Transformer\",\"activation\":\"gelu\",\"architectures\":[\"DistilBertModel\"],\"attention_dropout\":0.1,\"dim\":768,\"dropout\":0.1,\"hidden_dim\":3072,\"initializer_range\":0.02,\"max_position_embeddings\":512,\"model_type\":\"distilbert\",\"n_heads\":12,\"n_layers\":6,\"pad_token_id\":0,\"qa_dropout\":0.1,\"seq_classif_dropout\":0.2,\"sinusoidal_pos_embds\":false,\"tie_weights_\":true,\"transformers_version\":\"4.7.0\",\"vocab_size\":30522}"
  },
  "created_time": 1676073973126,
  "url": "https://artifacts.opensearch.org/models/ml-models/huggingface/sentence-transformers/msmarco-distilbert-base-tas-b/1.0.1/torch_script/sentence-transformers_msmarco-distilbert-base-tas-b-1.0.1-torch_script.zip"
}'


POST /_plugins/_ml/models/_register
{
  "name": "huggingface/sentence-transformers/msmarco-distilbert-base-tas-b",
  "version": "1.0.1",
  "model_group_id": "wlcnb4kBJ1eYAeTMHlV6",
  "description": "This is a port of the DistilBert TAS-B Model to sentence-transformers model: It maps sentences & paragraphs to a 768 dimensional dense vector space and is optimized for the task of semantic search.",
  "model_task_type": "TEXT_EMBEDDING",
  "model_format": "TORCH_SCRIPT",
  "model_content_size_in_bytes": 266352827,
  "model_content_hash_value": "acdc81b652b83121f914c5912ae27c0fca8fabf270e6f191ace6979a19830413",
  "model_config": {
    "model_type": "distilbert",
    "embedding_dimension": 768,
    "framework_type": "sentence_transformers",
    "all_config": """{"_name_or_path":"old_models/msmarco-distilbert-base-tas-b/0_Transformer","activation":"gelu","architectures":["DistilBertModel"],"attention_dropout":0.1,"dim":768,"dropout":0.1,"hidden_dim":3072,"initializer_range":0.02,"max_position_embeddings":512,"model_type":"distilbert","n_heads":12,"n_layers":6,"pad_token_id":0,"qa_dropout":0.1,"seq_classif_dropout":0.2,"sinusoidal_pos_embds":false,"tie_weights_":true,"transformers_version":"4.7.0","vocab_size":30522}"""
  },
  "created_time": 1676073973126,
  "url": "https://artifacts.opensearch.org/models/ml-models/huggingface/sentence-transformers/msmarco-distilbert-base-tas-b/1.0.1/torch_script/sentence-transformers_msmarco-distilbert-base-tas-b-1.0.1-torch_script.zip"
}
+ To check the status of the operation
  curl -XGET "http://localhost:9200/_plugins/_ml/tasks/cVeMb4kBJ1eYAeTMFFgj"

### Step 3: Deploy the model
curl -XPOST "http://localhost:9200/_plugins/_ml/models/cleMb4kBJ1eYAeTMFFg4/_deploy"

, check the status of the operation
curl -XGET "http://localhost:9200/_plugins/_ml/tasks/vVePb4kBJ1eYAeTM7ljG"

### Step 4 (Optional): Test the model
  curl -XPOST "http://localhost:9200/_plugins/_ml/_predict/text_embedding/cleMb4kBJ1eYAeTMFFg4" -H 'Content-Type: application/json' -d'
{
  "text_docs":[ "today is sunny"],
  "return_number": true,
  "target_response": ["sentence_embedding"]
}'
For a sparse encoding model

curl -XPOST "http://localhost:9200/_plugins/_ml/_predict/sparse_encoding/cleMb4kBJ1eYAeTMFFg4" -H 'Content-Type: application/json' -d'
{
  "text_docs":[ "today is sunny"]
}'

### Step 5: Use the model for search
  
  Question answering models

  - curl -XPOST "http://localhost:9200/_plugins/_ml/models/_register" -H 'Content-Type: application/json' -d'
{
    "name": "question_answering",
    "version": "1.0.0",
    "function_name": "QUESTION_ANSWERING",
    "description": "test model",
    "model_format": "TORCH_SCRIPT",
    "model_group_id": "lN4AP40BKolAMNtR4KJ5",
    "model_content_hash_value": "e837c8fc05fd58a6e2e8383b319257f9c3859dfb3edc89b26badfaf8a4405ff6",
    "model_config": { 
        "model_type": "bert",
        "framework_type": "huggingface_transformers"
    },
    "url": "https://github.com/opensearch-project/ml-commons/blob/main/ml-algorithms/src/test/resources/org/opensearch/ml/engine/algorithms/question_answering/question_answering_pt.zip?raw=true"
}'