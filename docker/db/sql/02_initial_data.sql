INSERT INTO `users` (`name`, `password`) VALUES ("admin", "$2a$10$sQ.xiNsXTZPCrn1bGRRqCuvkuP0RxxEiuGK6oCpfsp2BeoAfJTvPW");
INSERT INTO `categories` (`name`) VALUES ("sample-category-01");
INSERT INTO `tasks` (`title`, `user_id`, `category_id`) VALUES ("sample-task-01", 1, 1);
INSERT INTO `tasks` (`title`, `user_id`) VALUES ("sample-task-02", 1);
INSERT INTO `tasks` (`title`, `user_id`) VALUES ("sample-task-03", 1);
INSERT INTO `tasks` (`title`, `user_id`) VALUES ("sample-task-04", 1);
INSERT INTO `tasks` (`title`, `user_id`, `is_done`) VALUES ("sample-task-05", 1, true);
