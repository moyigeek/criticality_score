"use client";
import GitLinkTable from '@/components/GitLinkTable';
import TopNav from '@/components/TopNav';
import SearchBar from '@/components/SearchBar';
import SearchTable from '@/components/SearchTable';
import { Space } from 'antd';
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
            <Space direction="vertical" size="middle" style={{ display: 'flex' }}>
            <TopNav initialKey="gitlink"/>
            <SearchBar onSearch={handleSearch} onClearSearch={handleClearSearch} />
            {isSearching ? (
                <SearchTable searchQuery={searchQuery} />
            ) : (
                <GitLinkTable />
            )}
            </Space>
        </div>
    );
};

export default GitLinkPage;