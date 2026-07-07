# Team Service
This Is part of Submodule of [Warehouse Infra](https://github.com/pdcgo/warehouse_infra). In Warehouse Infra this is live in folder `./team_service`.<br>
This Service planned and Intended for replacing legacy Team System that exists in Warehouse Infra. Its planned for microservice and planned to more dependentless, separating domain purpose for better developing big and complex system that exists in Warehouse Infra.<br>
For now its just be candidate for Take over Team System In Warehouse Infra legacy.
Status for this development is still in progress and not completely take over legacy system.

1. for database schema related, read this [Database Schema](database-schema.md).
2. for testing that should cover and documentation about testing read this [Testing](test.md).

## Authentication & Authorization
1. Use v2 roling system. not legacy system. for complete reference read [this](../../user_service/docs/readme.md#authentication--authorization)
2. use interceptor that live in [here](../../user_service/access_interceptors/interceptor.go)
3. DON'T use legacy interface on [this](../../shared/interfaces/authorization_iface/authorization.go)


## Connect RPC Spec
`TeamService` heavyly depend `connect-rpc` to serve and creating apis and grpc. Why we use `connectrpc` because its can be two mode as pure grpc and grpc-web that interact like web. And also supported http2. This service have several rpc:

1. Team Management related RPC
    - Create Team that named `TeamCreate`
    - Delete Team that named `TeamDelete`
    - List Team that named `TeamList`
    - Update Team that named `TeamUpdate`
    - Detail Team that named `TeamDetail`

2. Team Info for preloading other service.
    - Getting Bulk Team named `TeamByIds`



### Team Management RPC
1. team management is for admin team.
2. `TeamList` and `TeamDetail` available for all authenticated user.
4. `TeamDelete` is soft delete.
5. `TeamCreate` is using implementation of [Implementation For Long Running Task RPC](../../docs/code-implementation-guideline.md#implementation-for-long-running-task-rpc).<br>
    Important step that execute in `TeamCreate` is:
    1. create team first.
    2. create team info.
    3. add who created as team owner.

## Team Info for preloading other service RPC
1. All authenticated user can access `TeamByIds`