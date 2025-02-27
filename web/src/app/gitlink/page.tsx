"use client";
import GitLinkTable from '@/components/GitLinkTable';
import TopNav from '@/components/TopNav';
import React from 'react';

const GitLinkPage: React.FC = () => {
    return (
        <div>
            <TopNav initialKey="gitlink"/>
            <GitLinkTable />
            
        </div>
    );
};

export default GitLinkPage;