import Nav from "@/components/Nav";
import React from "react";

type Props = {
  modal: React.ReactNode;
  results: React.ReactNode;
}

export default function ({
  modal, results
}: Props) {
  return <div>
    <Nav />
    {results}
    {modal}
  </div>
}