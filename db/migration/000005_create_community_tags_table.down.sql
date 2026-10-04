ALTER TABLE ONLY public.community_tags DROP CONSTRAINT community_tags_tag_id_fkey;
ALTER TABLE ONLY public.community_tags DROP CONSTRAINT community_tags_community_id_fkey;

DROP TABLE public.community_tags;