alter table public.scores_dist
    alter column score_id drop identity;

alter table public.scores_git
    alter column score_id drop identity;

alter table public.scores_lang
    alter column score_id drop identity;