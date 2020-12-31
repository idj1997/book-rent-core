insert into books values (10000, '1/1/2020', null, null, 'title1', 'content1', 5);
insert into books values (10001, '1/2/2020', null, null, 'title2', 'content2', 15);

-- password: 1234
insert into users values (10000, '1/1/2020', null, null, 'john', 'doe', 'johndoe@gmail.com', '$2y$12$Z51tvYyB2xEUejQydGcaiuCs1i3xqgHMvHwlzVLLQCk/7KVzahP9W', 0)
insert into users values (10001, '1/1/2020', null, null, 'mark', 'parker', 'markparker@gmail.com', '$2y$12$Z51tvYyB2xEUejQydGcaiuCs1i3xqgHMvHwlzVLLQCk/7KVzahP9W', 1)

insert into rent_details values (10000, '1/1/2020', null, null, 10000, 10000, 0)
insert into rent_details values (10001, '1/1/2020', null, null, 10001, 10000, 1)
insert into rent_details values (10002, '1/1/2020', null, null, 10000, 10001, 2)
insert into rent_details values (10003, '1/1/2020', null, null, 10001, 10001, 0)
