ALTER TABLE ONLY public.testimonials DROP CONSTRAINT testimonials_user_id_fkey;
ALTER TABLE ONLY public.testimonials DROP CONSTRAINT testimonials_pkey;

DROP TABLE public.testimonials;