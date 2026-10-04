CREATE TABLE public.community_members (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    community_id uuid NOT NULL,
    user_id uuid NOT NULL,
    role character varying DEFAULT 'Member'::character varying NOT NULL
);

ALTER TABLE ONLY public.community_members
    ADD CONSTRAINT community_members_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.community_members
    ADD CONSTRAINT community_members_community_id_fkey FOREIGN KEY (community_id) REFERENCES public.communities(id) DEFERRABLE;

ALTER TABLE ONLY public.community_members
    ADD CONSTRAINT community_members_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;