"use client";
import React, { useState } from "react";
import { Input, Button } from "antd";

interface SearchBarProps {
    onSearch: (query: string) => void;
    onClearSearch: () => void;
}

const SearchBar: React.FC<SearchBarProps> = ({ onSearch, onClearSearch }) => {
    const [searchQuery, setSearchQuery] = useState("");

    const handleSearch = () => {
        if (searchQuery.trim()) {
            onSearch(searchQuery);
        }
    };

    const handleClear = () => {
        setSearchQuery("");
        onClearSearch();
    };

    return (
        <div style={{ display: "flex", alignItems: "center" }}>
            <Input
                placeholder="Search packages"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                style={{ width: 300, marginRight: 10 }}
            />
            <Button type="primary" onClick={handleSearch} style={{ marginRight: 10 }}>
                Search
            </Button>
            <Button onClick={handleClear}>
                Clear
            </Button>
        </div>
    );
};

export default SearchBar;