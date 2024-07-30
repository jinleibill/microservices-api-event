# 商城类微服务

## 特性

- 事件驱动架构 - 事件溯源
- 六边形架构
- Grpc

## 运行
```
cd order
APPLICATION_PORT=8080 DATA_SOURCE_URL=192.168.64.7 go run main.go
```

## API
```
grpcurl -d '{"user_id": "123", "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' -plaintext localhost:8080 Order/Create 0
```

## 服务

- [订单服务](/order)

## TODO

- Saga
- BFF 服务网关