drop keyspace projectx;
create keyspace projectx
    with replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

create table projectx.message (
    uuid text primary key,
    title text,
    email text,
    content text,
    magic_number int,
    meta_created_tstamp timestamp
);

insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('d0a16893-55f5-4fbc-b88f-ada87215c363', 'interview 1', 'anna.zajkowska@example.com', 'simple text 1', 123, toTimestamp(now())) USING TTL 3000;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('c6e8b12a-ca23-43c9-aa6e-59e20e201c2c', 'interview 2', 'anna.zajkowska@example.com', 'simple text 2', 234, toTimestamp(now()))  USING TTL 3000;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 345, toTimestamp(now()))  USING TTL 3000;