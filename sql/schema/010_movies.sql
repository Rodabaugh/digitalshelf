-- +goose Up
alter table movies
add search tsvector
generated always as (
    setweight(to_tsvector('english',title), 'A')  || ' ' ||
    setweight(to_tsvector('simple',actors), 'B') || ' ' ||
    setweight(to_tsvector('english',genre), 'C') || ' ' ||
    setweight(to_tsvector('simple',writer), 'D') || ' ' ||
    setweight(to_tsvector('simple',director), 'D') :: tsvector
) stored;

create index idx_search on movies using GIN(search);
-- +goose Down
alter table movies
drop column search;
drop index idx_search;