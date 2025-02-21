create materialized view if not exists rankings as (
    select *, rank() over (order by score desc nulls last) as ranking
        from (select distinct on (ag.git_link) ag.git_link   as git_link,
                                                    s.id          as score_id,
                                                    s.dist_score  as dist_score,
                                                    s.lang_score  as lang_score,
                                                    s.git_score   as git_score,
                                                    s.score       as score,
                                                    s.update_time as update_time
            from all_gitlinks_cache ag
                    left join scores s on ag.git_link = s.git_link
            order by ag.git_link, s.id desc) as t
        order by score desc nulls last
);
