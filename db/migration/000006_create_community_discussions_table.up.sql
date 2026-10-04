CREATE TABLE public.community_discussions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    community_id uuid NOT NULL,
    user_id uuid NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.community_discussions
    ADD CONSTRAINT community_discussions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.community_discussions
    ADD CONSTRAINT community_discussions_community_id_fkey FOREIGN KEY (community_id) REFERENCES public.communities(id) DEFERRABLE;
  
ALTER TABLE ONLY public.community_discussions
    ADD CONSTRAINT community_discussions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;