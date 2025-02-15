"use client";
import DistroSel from '@/components/DistroSel';
import GitLinkTable from '@/components/GitLinkTable';
import React from 'react';

const GitLinkPage: React.FC = () => {
    return (
        <div>
            <DistroSel />
            <GitLinkTable />
        </div>
    );
};

export default GitLinkPage;