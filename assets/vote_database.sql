-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: Feb 28, 2025 at 03:12 PM
-- Server version: 10.4.28-MariaDB
-- PHP Version: 8.2.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `vote_database`
--

-- --------------------------------------------------------

--
-- Table structure for table `poll_info`
--

CREATE TABLE `poll_info` (
  `poll_id` int(20) NOT NULL,
  `poll_name` varchar(20) NOT NULL,
  `user_nickname` varchar(20) NOT NULL,
  `poll_question` varchar(50) NOT NULL,
  `poll_option1` varchar(30) NOT NULL,
  `poll_option2` varchar(30) NOT NULL,
  `user_vote_image` varchar(256) DEFAULT NULL COMMENT 'if user upload image'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `poll_info`
--

INSERT INTO `poll_info` (`poll_id`, `poll_name`, `user_nickname`, `poll_question`, `poll_option1`, `poll_option2`, `user_vote_image`) VALUES
(8, 'x', 'x', 'x', 'x', 'x', NULL),
(9, 'php', 'x', 'is php is good language?', 'yes', 'no', NULL),
(10, 'python', 'v', 'is python is good language?', 'yes', 'no', NULL),
(11, 'z', 'z', 'z', 'z', 'z', NULL),
(12, 'lol', 'z', 'lol', 'lol', 'lol', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `user_info`
--

CREATE TABLE `user_info` (
  `id` int(20) NOT NULL,
  `user_id` varchar(20) NOT NULL,
  `user_password` varchar(256) NOT NULL,
  `user_nickname` varchar(20) NOT NULL,
  `user_email` varchar(50) NOT NULL,
  `user_image` varchar(256) NOT NULL DEFAULT 'default_profile_image' COMMENT 'if user upload image'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user_info`
--

INSERT INTO `user_info` (`id`, `user_id`, `user_password`, `user_nickname`, `user_email`, `user_image`) VALUES
(10, 'x', '$2y$10$83axO.ts8ObIk6nnY5ZpqOjVabz1W7.NpWxz6OrCo5Ii6aAUNl4VW', 'x', 'x@x.x', 'default_profile_image'),
(11, 'v', '$2y$10$fPHTz464p1dViNWiTv1Fn.ebpgmNijHD1JY1Im/efZOJs1EUiynPe', 'v', 'v@v.v', 'default_profile_image'),
(12, 'z', '$2y$10$tL7MJnJi.0gQIWTFjqmgaOcl2wJjcn50N6OJki2xGbUXskIUq77o2', 'z', 'z@z.z', '/img_file/67c00efd2b560.png'),
(13, '=', '$2y$10$wAToVzyHmivZAGnlzU4pR.T9f7qg8tDp7VCyLRdaN4c1nklBBbBHO', '=', 'x@x.x', '/img_file/67c1c3c45290c.png'),
(14, '1', '$2y$10$FJUwy2nljoXlVRbPFwNcbetGYkyLOJgmLB32vIJn/vT3V.LvLILr.', '1', 'x@x.x', '/img_file/67c1c3e212b99.png');

-- --------------------------------------------------------

--
-- Table structure for table `vote_info`
--

CREATE TABLE `vote_info` (
  `vote_id` int(20) NOT NULL,
  `poll_id` int(20) NOT NULL,
  `user_nickname` varchar(20) NOT NULL,
  `vote_option` int(2) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `vote_info`
--

INSERT INTO `vote_info` (`vote_id`, `poll_id`, `user_nickname`, `vote_option`) VALUES
(10, 8, 'x', 1),
(11, 8, 'x', 1),
(12, 9, 'x', 1),
(13, 9, 'x', 1),
(14, 9, 'x', 2),
(15, 11, 'z', 1),
(16, 11, 'z', 2),
(17, 11, 'z', 1),
(18, 9, 'z', 2),
(19, 9, 'z', 2);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `poll_info`
--
ALTER TABLE `poll_info`
  ADD PRIMARY KEY (`poll_id`);

--
-- Indexes for table `user_info`
--
ALTER TABLE `user_info`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `vote_info`
--
ALTER TABLE `vote_info`
  ADD PRIMARY KEY (`vote_id`),
  ADD KEY `poll_id` (`poll_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `poll_info`
--
ALTER TABLE `poll_info`
  MODIFY `poll_id` int(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=13;

--
-- AUTO_INCREMENT for table `user_info`
--
ALTER TABLE `user_info`
  MODIFY `id` int(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT for table `vote_info`
--
ALTER TABLE `vote_info`
  MODIFY `vote_id` int(20) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=20;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `vote_info`
--
ALTER TABLE `vote_info`
  ADD CONSTRAINT `vote_info_ibfk_1` FOREIGN KEY (`poll_id`) REFERENCES `poll_info` (`poll_id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
