import React from "react";
import ReactDOM from "react-dom/client";

import './index.css'
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Start from './pages/Start.tsx';
import Profile from "./pages/Profile.tsx";
import Game from './pages/Game.tsx';
import Waiting from './pages/Waiting.tsx';
import QuestionSelector from "./pages/QuestionSelector.tsx";

const router = createBrowserRouter([
    {
        path: "/",
      element: <Start/>
    },
    {
        path: "/user",
        element: <Profile/>
    },
    {
        path: "/question",
        element: <QuestionSelector />,
    },
    {
        path:"/game/:id",
        element: <Game/>
    },
    {
        path: "/waiting",
        element: <Waiting/>
    }
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>
);
