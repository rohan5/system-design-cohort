# Kafka experiments

Setup Kafka locally to understand semantics and gurantees
MOST IMPORTANT Concepts

Partition
Offset
Consumer group
Ordering scope
Replayability
Rebalancing
Delivery semantics

### Kafka fundamentals learnt-
- Kafka guarantees ordering ONLY inside a partition. NOT globally.
- One partition can only be actively consumed by ONE consumer in a consumer group.
    - VERY IMPORTANT RULE -
        ```
        If:
        topic has 3 partitions
        Then:
        max active consumers in group = 3
        Extra consumers sit idle. 
        ```

- all old messages replay. This is Kafka’s superpower:
    - events persist
    - consumers track offsets independently
