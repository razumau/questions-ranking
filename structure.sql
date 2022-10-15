create table questions
(
	id serial not null
		constraint questions_pk
			primary key,
	question text,
	answer text
);

create table comparisons
(
	id uuid default gen_random_uuid() not null
		constraint comparisons_pk
			primary key,
	first integer
		constraint first
			references questions,
	second integer
		constraint second
			references questions,
	selected integer
		constraint selected
			references questions
);
