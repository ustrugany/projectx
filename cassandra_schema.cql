-- drop keyspace projectx;
create keyspace projectx
    with replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};

-- drop table projectx.message;
create table projectx.message (
    uuid text,
    title text,
    email text,
    content text,
    magic_number int,
    meta_created_tstamp timestamp,
    // assuming that magic number will have higher cardinality than email
    primary key ((magic_number, email))
);

insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('d0a16893-55f5-4fbc-b88f-ada87215c363', 'interview 1', 'anna.zajkowska@example.com', 'simple text 1', 123, toTimeStamp(now())) USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('c6e8b12a-ca23-43c9-aa6e-59e20e201c2c', 'interview 2', 'anna.zajkowska@example.com', 'simple text 2', 234, toTimestamp(now()))  USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 345, toTimestamp(now()))  USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 456, toTimestamp(now()))  USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 567, toTimestamp(now()))  USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 678, toTimestamp(now()))  USING TTL 600;
insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp)
values ('9d86c8a3-3835-4aa1-b016-495bf5c469a7', 'interview 3', 'anna.zajkowska@example.com', 'simple text 3', 890, toTimestamp(now()))  USING TTL 600;
