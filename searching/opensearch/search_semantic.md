
# PREREQUISITE
    - Choosing a model
# Process 
## Step 1: Create an ingest pipeline

curl -k -XPUT "https://localhost:9200/_ingest/pipeline/nlp-ingest-pipeline" -H 'Content-Type: application/json' -d'
{
  "description": "A text embedding pipeline",
  "processors": [
    {
      "text_embedding": {
        "model_id": "xokaEJABlLB-9JtUZovQ",
        "field_map": {
          "passage_text": "passage_embedding"
        }
      }
    }
  ]
}' -u admin:admin 

+ Delete:
   curl -k -XDELETE "https://localhost:9200/_ingest/pipeline/nlp-ingest-pipeline" -u admin:admin
## Step 2: Create an index for ingestion
  curl -k   -XPUT "https://localhost:9200/my-nlp-index" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "index.knn": true,
    "default_pipeline": "nlp-ingest-pipeline"
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "text"
      },
      "passage_embedding": {
        "type": "knn_vector",
        "dimension": 768,
        "method": {
          "engine": "lucene",
          "space_type": "l2",
          "name": "hnsw",
          "parameters": {}
        }
      },
      "passage_text": {
        "type": "text"
      }
    }
  }
}' -u admin:admin 

+ Delete:
  curl -k -XDELETE "https://localhost:9200/my-nlp-index" -u admin:admin 


## Step 3: Ingest documents into the index
  curl -k -XPUT "https://localhost:9200/my-nlp-index/_doc/1" -H 'Content-Type: application/json' -d'
{
  "passage_text": "Hello world",
  "id": "s1"
}' -u admin:admin  

curl -k -XPUT "https://localhost:9200/my-nlp-index/_doc/2" -H 'Content-Type: application/json' -d'
{
  "passage_text": "Hi planet",
  "id": "s2"
}' -u admin:admin  

## Step 4: Search the index using neural search

curl -k -XGET "https://localhost:9200/my-nlp-index/_search" -H 'Content-Type: application/json' -d'
{
  "_source": {
    "excludes": [
      "passage_embedding"
    ]
  },
  "query": {
    "bool": {
      "filter": {
         "wildcard":  { "id": "*1" }
      },
      "should": [
        {
          "script_score": {
            "query": {
              "neural": {
                "passage_embedding": {
                  "query_text": "Hi world",
                  "model_id": "xokaEJABlLB-9JtUZovQ",
                  "k": 100
                }
              }
            },
            "script": {
              "source": "_score * 1.5"
            }
          }
        },
        {
          "script_score": {
            "query": {
              "match": {
                "passage_text": "Hi world"
              }
            },
            "script": {
              "source": "_score * 1.7"
            }
          }
        }
      ]
    }
  }
}' -u admin:admin 

# Setting a default model on an index or field
- curl -k -XPUT "https://localhost:9200/_search/pipeline/default_model_pipeline" -H 'Content-Type: application/json' -d'
{
  "request_processors": [
    {
      "neural_query_enricher" : {
        "default_model_id": "xokaEJABlLB-9JtUZovQ",
        "neural_field_default_id": {
           "my_field_1": "xokaEJABlLB-9JtUZovQ",
           "my_field_2": "xokaEJABlLB-9JtUZovQ"
        }
      }
    }
  ]
}' -u admin:admin

# Then set the default model for your index:

curl -XPUT "http://localhost:9200/my-nlp-index/_settings" -H 'Content-Type: application/json' -d'
{
  "index.search.default_pipeline" : "default_model_pipeline"
}'


#  omit the model ID when searching:
  curl -XGET "http://localhost:9200/my-nlp-index/_search" -H 'Content-Type: application/json' -d'
{
  "_source": {
    "excludes": [
      "passage_embedding"
    ]
  },
  "query": {
    "neural": {
      "passage_embedding": {
        "query_text": "Hi world",
        "k": 100
      }
    }
  }
}'
