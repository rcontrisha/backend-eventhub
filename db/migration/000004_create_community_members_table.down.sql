ALTER TABLE ONLY public.community_members DROP CONSTRAINT community_members_user_id_fkey;
ALTER TABLE ONLY public.community_members DROP CONSTRAINT community_members_community_id_fkey;
ALTER TABLE ONLY public.community_members DROP CONSTRAINT community_members_pkey;

DROP TABLE public.community_members;