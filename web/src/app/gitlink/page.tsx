"use client";
import GitLinkTable from '@/components/GitLinkTable';
import TopNav from '@/components/TopNav';
import SearchBar from '@/components/SearchBar';
import SearchTable from '@/components/SearchTable';
import React, { useState } from 'react';

const GitLinkPage: React.FC = () => {
    const [searchQuery, setSearchQuery] = useState("");
    const [isSearching, setIsSearching] = useState(false);

    const handleSearch = (query: string) => {
        setSearchQuery(query);
        setIsSearching(true);
    };

    const handleClearSearch = () => {
        setSearchQuery("");
        setIsSearching(false);
    };

    return (
        <div>
            <TopNav initialKey="gitlink"/>
            <SearchBar onSearch={handleSearch} onClearSearch={handleClearSearch} />
            {isSearching ? (
                <SearchTable searchQuery={searchQuery} />
            ) : (
                <GitLinkTable />
            )}
        </div>
    );
};

export default GitLinkPage;