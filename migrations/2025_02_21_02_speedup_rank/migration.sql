drop materialized view if exists rankings;
create or replace view rankings as (
    select *, rank() over (order by score desc nulls last) as ranking
            from (select  s.git_link   as git_link,
                        s.id          as score_id,
                        s.dist_score  as dist_score,
                        s.lang_score  as lang_score,
                        s.git_score   as git_score,
                        s.score       as score,
                        s.update_time as update_time
                        from scores s
                where s.round = (select max(round) from scores)) as t
    order by score desc nulls last
);
create table rankings_cache as select * from rankings;