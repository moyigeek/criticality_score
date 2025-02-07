create table if not exists all_gitlinks_cache as select * from all_gitlinks;

create index if not exists scores_git_link_index
    on scores (git_link);
