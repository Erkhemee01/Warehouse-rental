import React, { useContext, Suspense, useEffect, lazy } from "react";
import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import routes from "../routes";

import Sidebar from "../components/Sidebar";
import Header from "../components/Header";
import Main from "./Main";
import ThemedSuspense from "../components/ThemedSuspense";
import { SidebarContext } from "../context/SidebarContext";

import Page404 from "../pages/404";

function Layout() {
  const { isSidebarOpen, closeSidebar } = useContext(SidebarContext);
  let location = useLocation();

  useEffect(() => {
    closeSidebar();
  }, [location]);

  return (
    <div
      className={`flex h-screen bg-gray-50 dark:bg-gray-900 ${
        isSidebarOpen && "overflow-hidden"
      }`}
    >
      <Sidebar />

      <div className="flex flex-col flex-1 w-full">
        <Header />
        <Main>
          <Suspense fallback={<ThemedSuspense />}>
            <Routes>
              {routes.map((route, i) => {
                return route.component ? (
                  <Route
                    key={i}
                    exact={true}
                    path={`/app${route.path}`}
                    render={(props) => <route.component {...props} />}
                  />
                ) : null;
              })}
              <Navigate exact from="/app" to="/app/dashboard" />
              <Route component={Page404} />
            </Routes>
          </Suspense>
        </Main>
      </div>
    </div>
  );
}

export default Layout;
