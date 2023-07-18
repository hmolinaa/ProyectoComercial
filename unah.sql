-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Servidor: 127.0.0.1
-- Tiempo de generación: 18-07-2023 a las 04:25:29
-- Versión del servidor: 10.4.28-MariaDB
-- Versión de PHP: 8.2.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de datos: `unah`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `students`
--

CREATE TABLE `students` (
  `id` int(5) NOT NULL,
  `name` varchar(25) NOT NULL,
  `account` int(5) NOT NULL,
  `subject` varchar(10) NOT NULL,
  `first_partial` int(2) NOT NULL,
  `second_partial` int(2) NOT NULL,
  `third_partial` int(2) NOT NULL,
  `final_score` int(5) NOT NULL,
  `email` varchar(2500) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- RELACIONES PARA LA TABLA `students`:
--

--
-- Volcado de datos para la tabla `students`
--

INSERT INTO `students` (`id`, `name`, `account`, `subject`, `first_partial`, `second_partial`, `third_partial`, `final_score`, `email`) VALUES
(1, 'Nahun', 202015, 'MM-520', 78, 85, 88, 84, 'nahun.mart@gmail.com'),
(2, 'Ivan', 20356, 'MM-520', 70, 70, 70, 70, 'ivan@gmail.com'),
(3, 'Ana Gomez', 2019266, 'MM-520', 85, 80, 63, 84, 'ana.gomez@gmail.com'),
(10, 'Molina', 20181510, 'MM-520', 95, 95, 95, 95, 'hmolina100@gmail.com'),
(38, 'Henrry', 2147483, 'MM-520', 90, 90, 90, 90, 'hmolinaa@unah.hn'),
(45, 'Pedro', 25544, 'MM-520', 80, 8, 88, 90, 'pedro@gmail.com');

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `students`
--
ALTER TABLE `students`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `students`
--
ALTER TABLE `students`
  MODIFY `id` int(5) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=47;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
