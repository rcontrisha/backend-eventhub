CREATE TABLE public.community_tags (
    community_id uuid NOT NULL,
    tag_id uuid NOT NULL
);

ALTER TABLE ONLY public.community_tags
    ADD CONSTRAINT community_tags_community_id_fkey FOREIGN KEY (community_id) REFERENCES public.communities(id) DEFERRABLE;

ALTER TABLE ONLY public.community_tags
    ADD CONSTRAINT community_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) DEFERRABLE;